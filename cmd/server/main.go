package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sadia-54/QStack-Backend/internal/config"
)

func main() {
	env := config.Load()
	config.ConnectDB(env)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	e.Logger.Fatal(e.Start(":" + env.AppPort))
}