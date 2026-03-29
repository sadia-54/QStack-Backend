package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/api/handlers"
)

func RegisterUploadRoutes(api *echo.Group) {

	uploadHandler := handlers.NewUploadHandler()

	uploads := api.Group("/upload")

	uploads.POST("", uploadHandler.UploadImage)
}