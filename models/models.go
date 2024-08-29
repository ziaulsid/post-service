package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID                uint64         `json:"id"`
	Caption           sql.NullString `json:"caption"`
	OriginalImage     sql.NullString `json:"original_image"`
	OriginalImageName sql.NullString `json:"original_image_name"`
	ResizedImage      sql.NullString `json:"resized_image"`
	CreatedBy         uint64         `json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         sql.NullTime   `json:"deleted_at"`
}

type Comment struct {
	ID        uint64       `json:"id"`
	Content   string       `json:"content"`
	CreatedBy uint64       `json:"created_by"`
	PostID    uint64       `json:"post_id"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type User struct {
	ID        uint64       `json:"id"`
	Username  string       `json:"username"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type PostDetail struct {
	Post     *Post
	Comments []*Comment
}
