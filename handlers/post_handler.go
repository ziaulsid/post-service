package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"post-service/dto"
	"post-service/services"
	"strconv"
	"strings"
)

const maxImageUploadSize = 100000000

var (
	supportedFileExtensions = []string{"jpg", "jpeg", "png", "bmp"}
)

type PostHandler struct {
	postService services.PostService
}

func NewPostHandler(postService services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePostHandler(c *gin.Context) {
	var req dto.CreatePostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.UserID <= 0 {
		log.Error().Msg("Invalid User ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	post, err := h.postService.CreatePost(req.Caption, req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create post")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, dto.CreatePostResponse{
		ID:      post.ID,
		Caption: post.Caption.String,
		UserID:  post.CreatedBy,
	})
}

func (h *PostHandler) GetPostsByUser(c *gin.Context) {
	userIdStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	posts, err := h.postService.GetPostsByUser(userID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting posts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) AddImageToPost(context *gin.Context) {
	postIDStr := context.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing post id")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing the post_id"})
		return
	}

	file, metadata, err := context.Request.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("Error parsing file image")
		context.JSON(http.StatusBadRequest, gin.H{"error": "No image provided for the post"})
		return
	}

	if metadata.Size > maxImageUploadSize {
		log.Error().Err(err).Msg("File uploads greater than 100 MB are not supported")
		context.JSON(http.StatusBadRequest, gin.H{"error": "File uploads greater than 100 MB are not supported"})
		return
	}

	if !h.isFileExtSupported(metadata.Filename) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "File extension is not supported"})
		return
	}

	err = h.postService.AttachImage(postID, file, metadata)
	if err != nil {
		log.Error().Err(err).Msg("Error adding image to post")
		context.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Server error while fetching post: %d", postID)})
		return
	}

	context.JSON(http.StatusNoContent, gin.H{})
}

func (h *PostHandler) isFileExtSupported(filename string) bool {
	for _, extension := range supportedFileExtensions {
		if strings.Contains(filename, fmt.Sprintf(".%s", extension)) {
			return true
		}
	}

	return false
}
