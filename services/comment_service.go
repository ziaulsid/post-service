package services

import (
	"post-service/dto"
	"post-service/models"
	"post-service/repositories"
)

type CommentService interface {
	AddComment(postID, userId uint64, content string) (*dto.AddCommentResponse, error)
	DeleteComment(id string) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
}

func NewCommentService(commentRepo repositories.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) AddComment(postID, userId uint64, content string) (*dto.AddCommentResponse, error) {
	comment := &models.Comment{
		PostID:    postID,
		Content:   content,
		CreatedBy: userId,
	}

	err := s.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return &dto.AddCommentResponse{
		PostID:  postID,
		Content: content,
		UserID:  userId,
	}, nil

}

func (s *commentService) DeleteComment(id string) error {
	return s.commentRepo.DeleteByID(id)
}
