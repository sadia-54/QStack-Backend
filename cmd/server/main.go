package main

import (
	"net/http"
	"log"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/sadia-54/qstack-backend/internal/config"
	customMiddleware "github.com/sadia-54/qstack-backend/internal/api/middleware"
	"github.com/sadia-54/qstack-backend/internal/api/routes"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
	"github.com/sadia-54/qstack-backend/internal/validator"
	"github.com/sadia-54/qstack-backend/internal/services"
	"github.com/sadia-54/qstack-backend/internal/repositories"
	"github.com/sadia-54/qstack-backend/internal/queue"
)

func main() {
	env := config.Load() // load env
	config.ConnectDB(env) // connect to DB

	// connect to RabbitMQ
	if err := queue.Connect(); err != nil {
		log.Fatal(err)
	}
	defer queue.Close()

	// Initialize repositories
	// user repo
	userRepo := repositories.NewUserRepository(config.DB)
	tokenRepo := repositories.NewEmailVerificationTokenRepository(config.DB)

	// question repo
	questionRepo := repositories.NewQuestionRepository(config.DB)
	tagRepo := repositories.NewTagRepository(config.DB)
	voteRepo := repositories.NewQuestionVoteRepository(config.DB)

	// answer repo
	answerRepo := repositories.NewAnswerRepository(config.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo, env.JWTSecret, env.AppBaseURL)
	questionService := services.NewQuestionService(questionRepo, tagRepo, voteRepo)
	answerService := services.NewAnswerService(answerRepo, questionRepo)
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	answerHandler := handlers.NewAnswerHandler(answerService)
	userHandler := handlers.NewUserHandler(userService)

	// setup echo server
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// serve uploaded images
	e.Static("/uploads", "uploads")

	// CORS
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

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
	// register question routes
	routes.RegisterQuestionRoutes(api, questionHandler)
	// register answer routes
	routes.RegisterAnswerRoutes(api, answerHandler)
	// register user routes
	routes.RegisterUserRoutes(api, userHandler)

	// image upload routes
	routes.RegisterUploadRoutes(api)

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