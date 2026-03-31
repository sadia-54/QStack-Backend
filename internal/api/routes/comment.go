package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/api/middleware"
)

func RegisterCommentRoutes(api *echo.Group, handler *handlers.CommentHandler) {

	comments := api.Group("/comments")

	protected := comments.Group("")
	protected.Use(middleware.JWTMiddleware())

	protected.POST("/answer/:answer_id", handler.Create)
	protected.PUT("/:id", handler.Update)    
	protected.DELETE("/:id", handler.Delete)

	comments.GET("/answer/:answer_id", handler.GetByAnswer)
}