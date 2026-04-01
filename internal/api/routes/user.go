package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/api/middleware"
)

func RegisterUserRoutes(e *echo.Group, userHandler *handlers.UserHandler) {

	user := e.Group("/users")
	user.GET("/community/stats", userHandler.GetCommunityStats)

	// Public routes
	user.GET("/:id/profile", userHandler.GetProfile)
	user.GET("", userHandler.GetUsers)

	// Protected routes
	protected := user.Group("")
	protected.Use(middleware.JWTMiddleware())

	protected.PUT("/profile", userHandler.UpdateProfile)
	protected.GET("/:id/activity", userHandler.GetActivity)
	protected.GET("/me", userHandler.GetMyProfile)
	protected.POST("/profile/image", userHandler.UploadProfileImage)
}