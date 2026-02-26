package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/sadia-54/qstack-backend/internal/api/handlers"
)

func RegisterAuthRoutes(e *echo.Group, authHandler *handlers.AuthHandler) {
	auth := e.Group("/auth")

	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)
	auth.GET("/verify-email", authHandler.VerifyEmail)
}