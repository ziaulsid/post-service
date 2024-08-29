package repositories

import (
	"database/sql"
	"fmt"
	"post-service/models"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	DeleteByID(id string) error
	GetLastTwoComments(postID uint64) ([]*models.Comment, error)
}

type commentRepositoryStruct struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryStruct{db: db}
}

func (r *commentRepositoryStruct) Create(comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, content, created_by)
		VALUES (?, ?, ?)
	`

	_, err := r.db.Exec(query, comment.PostID, comment.Content, comment.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to create comment: %v", err)
	}

	return nil
}

func (r *commentRepositoryStruct) DeleteByID(id string) error {
	query := `
		DELETE FROM comments WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
	}

	return nil
}

func (r *commentRepositoryStruct) GetLastTwoComments(postID uint64) ([]*models.Comment, error) {
	query := `SELECT id, content, created_by, created_at 
              FROM comments WHERE post_id = ? ORDER BY created_at DESC LIMIT 2`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		if err = rows.Scan(&comment.ID, &comment.Content, &comment.CreatedBy, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
