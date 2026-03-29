package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c echo.Context) error {

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "image file required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "cannot read file",
		})
	}
	defer src.Close()

	// create uploads folder if not exists
	os.MkdirAll("uploads", os.ModePerm)

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	path := filepath.Join("uploads", filename)

	dst, err := os.Create(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "cannot save file",
		})
	}
	defer dst.Close()

	if _, err = dst.ReadFrom(src); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to write file",
		})
	}

	url := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	return c.JSON(http.StatusOK, echo.Map{
		"url": url,
	})
}