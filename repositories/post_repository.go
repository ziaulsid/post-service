package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"post-service/models"
	"strings"
)

type PostRepository interface {
	Create(post *models.Post) (uint64, error)
	GetByID(postID uint64) (*models.Post, error)
	GetPostsByUserID(userID uint64) ([]*models.Post, error)
	Update(post *models.Post, updateFields map[string]interface{}) (*models.Post, error)
}

type postRepositoryStruct struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepositoryStruct{db: db}
}

func (r *postRepositoryStruct) Create(post *models.Post) (uint64, error) {
	query := `
		INSERT INTO posts (caption, created_by)
		VALUES (?, ?)
	`

	result, err := r.db.Exec(query, post.Caption, post.CreatedBy)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}
	return uint64(id), nil
}

func (r *postRepositoryStruct) GetByID(postID uint64) (*models.Post, error) {
	query := `
			SELECT id, 
       				caption, 
       				resized_image, 
       				created_by, 
       				created_at 
            FROM posts 
            WHERE id = ?
`

	post := &models.Post{}
	err := r.db.QueryRow(query, postID).Scan(
		&post.ID, &post.Caption, &post.ResizedImage, &post.CreatedBy, &post.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}

func (r *postRepositoryStruct) Update(post *models.Post, updateFields map[string]interface{}) (*models.Post, error) {
	setClauses := []string{}
	args := []interface{}{}

	for field, value := range updateFields {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}

	args = append(args, post.ID)
	query := fmt.Sprintf("UPDATE posts SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepositoryStruct) GetPostsByUserID(userID uint64) ([]*models.Post, error) {
	query := `SELECT id, caption, resized_image, created_at 
              FROM posts WHERE created_by = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(&post.ID, &post.Caption, &post.ResizedImage, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
