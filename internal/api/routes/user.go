package routes

import (

	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
)

func RegisterUserRoutes(e *echo.Group, userHandler *handlers.UserHandler) {
	user := e.Group("/users")
	user.GET("/:id/profile", userHandler.GetProfile)
	user.PUT("/profile", userHandler.UpdateProfile)
}