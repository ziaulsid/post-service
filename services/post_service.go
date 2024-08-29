package services

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
	"post-service/dto"
	"post-service/models"
	"post-service/repositories"
	"post-service/workers"
)

type PostService interface {
	CreatePost(caption string, userID uint64) (*models.Post, error)
	GetPostsByUser(userID uint64) (dto.PostsResponse, error)
	AttachImage(postID uint64, file multipart.File, metadata *multipart.FileHeader) error
}

type postService struct {
	postRepo    repositories.PostRepository
	commentRepo repositories.CommentRepository
	workerPool  *workers.WorkerPool
}

func NewPostService(
	postRepo repositories.PostRepository,
	commentRepo repositories.CommentRepository,
	workerPool *workers.WorkerPool,
) PostService {
	return &postService{postRepo: postRepo, commentRepo: commentRepo, workerPool: workerPool}
}

func (s *postService) AttachImage(postID uint64, file multipart.File, metadata *multipart.FileHeader) error {

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		log.Error().Err(err).Msg("Error occurred while fetching post from DB")
		return err
	}

	go func(post *models.Post) {
		defer func() {
			_ = file.Close()
		}()

		dst, err := os.Create(metadata.Filename)
		if err != nil {
			log.Error().Err(err).Msg("Error occurred while creating file")
			return
		}

		defer func() {
			_ = dst.Close()
		}()
		if _, err = io.Copy(dst, file); err != nil {
			log.Error().Err(err).Msg("Error occurred while copying file")
			return
		}

		fInfo, _ := dst.Stat()
		size := fInfo.Size()
		buf := make([]byte, size)

		fReader := bufio.NewReader(dst)
		_, _ = fReader.Read(buf)

		imgBase64Str := base64.StdEncoding.EncodeToString(buf)

		post.OriginalImage = sql.NullString{
			String: imgBase64Str,
			Valid:  true,
		}

		post.OriginalImageName = sql.NullString{
			String: metadata.Filename,
			Valid:  true,
		}

		updateFields := map[string]interface{}{
			"original_image":      post.OriginalImage,
			"original_image_name": post.OriginalImageName,
		}

		_, err = s.postRepo.Update(post, updateFields)
		if err != nil {
			log.Printf("Error updating post with original image: %v", err)
		}

		s.workerPool.AddJob(workers.ImageJob{
			Post:      post,
			ImagePath: metadata.Filename,
			Width:     600,
			Height:    600,
		})

	}(post)

	return nil
}

func (s *postService) CreatePost(caption string, userID uint64) (*models.Post, error) {
	post := &models.Post{
		Caption: sql.NullString{
			String: caption,
			Valid:  true,
		},
		CreatedBy: userID,
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) GetPostsByUser(userID uint64) (dto.PostsResponse, error) {
	posts, err := s.postRepo.GetPostsByUserID(userID)
	if err != nil {
		return dto.PostsResponse{}, err
	}

	// For each post, fetch the last 2 comments and convert to response structs
	var postResponses []dto.PostResponse
	for _, post := range posts {
		comments, err := s.commentRepo.GetLastTwoComments(post.ID)
		if err != nil {
			return dto.PostsResponse{}, err
		}

		// Convert comments to response structs
		commentResponses := make([]dto.CommentResponse, len(comments))
		for i, comment := range comments {
			commentResponses[i] = dto.CommentResponse{
				ID:        comment.ID,
				Content:   &comment.Content,
				UserID:    comment.CreatedBy,
				CreatedAt: &comment.CreatedAt,
			}
		}

		var resizedImage string
		if post.ResizedImage.Valid {
			resizedImage = post.ResizedImage.String
		}

		// Convert post to response struct
		postResponse := dto.PostResponse{
			ID:           post.ID,
			Caption:      nullableString(post.Caption),
			ResizedImage: &resizedImage,
			CreatedAt:    &post.CreatedAt,
			Comments:     commentResponses,
		}

		postResponses = append(postResponses, postResponse)
	}

	return dto.PostsResponse{Posts: postResponses}, nil
}

func nullableString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
