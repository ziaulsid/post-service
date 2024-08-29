package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"post-service/dto"
	"post-service/services"
	"strconv"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) AddCommentHandler(c *gin.Context) {
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing postId")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing the post_id"})
		return
	}

	var req dto.AddCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to parse AddCommentRequest")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.UserID <= 0 {
		log.Error().Msg("UserId is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	if len(req.Content) <= 0 {
		log.Error().Msg("Content is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty comment"})
		return
	}

	response, err := h.commentService.AddComment(postID, req.UserID, req.Content)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add comment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) DeleteCommentHandler(c *gin.Context) {

}
