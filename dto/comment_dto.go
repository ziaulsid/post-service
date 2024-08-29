package dto

type AddCommentRequest struct {
	Content string `json:"content"`
	UserID  uint64 `json:"user_id"`
}

type AddCommentResponse struct {
	PostID  uint64 `json:"post_id"`
	Content string `json:"content"`
	UserID  uint64 `json:"creator"`
}
