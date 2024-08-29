package dto

import "time"

type CreatePostRequest struct {
	Caption string `json:"caption"`
	UserID  uint64 `json:"user_id"`
}

type PostResponse struct {
	ID           uint64            `json:"id"`
	Caption      *string           `json:"caption,omitempty"`
	ResizedImage *string           `json:"resized_image,omitempty"`
	CreatedAt    *time.Time        `json:"created_at,omitempty"`
	Comments     []CommentResponse `json:"comments,omitempty"`
}

type CommentResponse struct {
	ID        uint64     `json:"id"`
	UserID    uint64     `json:"user_id"`
	Content   *string    `json:"content,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type PostsResponse struct {
	Posts []PostResponse `json:"posts,omitempty"`
}
