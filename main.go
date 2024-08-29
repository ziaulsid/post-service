package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"post-service/handlers"
	"post-service/internal/db"
	internalLogger "post-service/internal/logger"
	"post-service/repositories"
	"post-service/services"
	"post-service/workers"
)

func main() {
	internalLogger.InitLogger()

	internalDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer internalDB.Close()

	postRepo := repositories.NewPostRepository(internalDB)
	commentRepo := repositories.NewCommentRepository(internalDB)

	wp := workers.NewWorkerPool(4, postRepo)

	postService := services.NewPostService(postRepo, commentRepo, wp)
	commentService := services.NewCommentService(commentRepo)

	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	r := gin.Default()

	// Post routes
	r.POST("/posts", postHandler.CreatePostHandler)
	r.GET("/users/:userId/posts", postHandler.GetPostsByUser)
	r.POST("/posts/:postId/images", postHandler.AddImageToPost)

	// Comment routes
	r.POST("/posts/:postId/comments", commentHandler.AddCommentHandler)
	r.DELETE("/posts/:postId/comments/:commentID", commentHandler.DeleteCommentHandler)

	// Start the server
	r.Run(":8080")

}
