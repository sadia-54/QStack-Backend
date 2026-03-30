package routes

import (
	"github.com/sadia-54/qstack-backend/internal/api/middleware"
	"github.com/labstack/echo/v4"

	"github.com/sadia-54/qstack-backend/internal/api/handlers"
)

func RegisterAuthRoutes(e *echo.Group, authHandler *handlers.AuthHandler) {
	auth := e.Group("/auth")
	protected := auth.Group("")
	protected.Use(middleware.JWTMiddleware())

	// Public routes
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)
	auth.GET("/verify-email", authHandler.VerifyEmail)

	// Protected routes
	protected.POST("/change-password", authHandler.ChangePassword)
}