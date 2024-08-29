package workers

import (
	"database/sql"
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"post-service/imageprocessor"
	"post-service/models"
	"post-service/repositories"
	"sync"
)

type WorkerPool struct {
	Jobs     chan ImageJob
	Wg       sync.WaitGroup
	postRepo repositories.PostRepository
}

type ImageJob struct {
	Post      *models.Post
	ImagePath string
	Width     uint
	Height    uint
}

func NewWorkerPool(numWorkers int, postRepo repositories.PostRepository) *WorkerPool {
	wp := &WorkerPool{
		Jobs:     make(chan ImageJob, 100),
		postRepo: postRepo,
	}

	wp.Wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go wp.worker(i)
	}

	return wp
}

func (wp *WorkerPool) worker(id int) {
	defer wp.Wg.Done()
	for job := range wp.Jobs {
		log.Printf("Worker %d: Processing image for post ID %d\n", id, job.Post.ID)

		resizedImage, err := imageprocessor.ResizeImage(job.ImagePath, job.Width, job.Height)
		if err != nil {
			log.Printf("Worker %d: Failed to resize image: %v\n", id, err)
			continue
		}

		updateFields := map[string]interface{}{
			"resized_image": sql.NullString{
				String: base64.StdEncoding.EncodeToString(resizedImage),
				Valid:  true,
			},
		}

		_, err = wp.postRepo.Update(job.Post, updateFields)
		if err != nil {
			log.Printf("Worker %d: Failed to update post with resized image: %v\n", id, err)
			continue
		}

		log.Printf("Worker %d: Finished processing image for post ID %d\n", id, job.Post.ID)
	}
}

func (wp *WorkerPool) AddJob(job ImageJob) {
	wp.Jobs <- job
}

func (wp *WorkerPool) Stop() {
	close(wp.Jobs)
	wp.Wg.Wait()
}
