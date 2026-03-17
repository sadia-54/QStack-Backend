package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/api/middleware"
)

func RegisterAnswerRoutes(api *echo.Group, handler *handlers.AnswerHandler) {

	answers := api.Group("/answers")

	protected := answers.Group("")
	protected.Use(middleware.JWTMiddleware())

	protected.POST("/question/:question_id", handler.Create)
	protected.PUT("/:id", handler.Update)
	protected.DELETE("/:id", handler.Delete)
	protected.PUT("/:id/accept", handler.Accept)

	answers.GET("/question/:question_id", handler.GetByQuestion)
}