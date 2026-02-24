package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/sadia-54/qstack-backend/internal/config"
	customMiddleware "github.com/sadia-54/qstack-backend/internal/api/middleware"
	"github.com/sadia-54/qstack-backend/internal/api/routes"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/validator"
	"github.com/sadia-54/qstack-backend/internal/services"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

func main() {
	env := config.Load() // load env
	config.ConnectDB(env) // connect to DB

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.DB)
	tokenRepo := repositories.NewEmailVerificationTokenRepository(config.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, env.JWTSecret, env.AppBaseURL)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// setup echo server
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// register validator
	e.Validator = validators.NewValidator()

	// server health check route
	e.GET("/health", func(c echo.Context) error {
		sqlDB, err := config.DB.DB()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "db handle error",
			})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "db not reachable",
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"db":     "connected",
		})
	})

	// api routes
	api := e.Group("/api/v1")

	// register auth routes
	routes.RegisterAuthRoutes(api, authHandler)

	// protected routes
	protected := api.Group("/protected")
	protected.Use(customMiddleware.JWTMiddleware())

	protected.GET("/me", func(c echo.Context) error {
		userID := c.Get("user_id")
		return c.JSON(http.StatusOK, echo.Map{
			"message": "This is a protected route",
			"user_id": userID,
		})
	})

	// start the server
	e.Logger.Fatal(e.Start(":" + env.AppPort))
}