package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/api/middleware"
)

func RegisterQuestionRoutes(api *echo.Group, handler *handlers.QuestionHandler) {

	questions := api.Group("/questions")

	// Public
	questions.GET("", handler.Feed)
	questions.GET("/:id", handler.GetByID)

	// Protected
	protected := questions.Group("")
	protected.Use(middleware.JWTMiddleware())

	protected.POST("", handler.Create)
	protected.PUT("/:id", handler.Update)
	protected.DELETE("/:id", handler.Delete)
	protected.GET("/my", handler.MyQuestions)

	// for voting
	protected.POST("/:id/vote", handler.Vote)

	// for user-specific feed
	protected.GET("/my-feed", handler.MyFeed)
}