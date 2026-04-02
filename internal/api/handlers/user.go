package handlers

import (
	"io"
	"os"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetProfile(c echo.Context) error {

	userIDParam := c.Param("id")
	userID, _ := strconv.ParseInt(userIDParam, 10, 64)

	res, err := h.userService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {

	userID := int64(c.Get("user_id").(float64))

	var req struct {
		Bio string `json:"bio"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.userService.UpdateProfile(userID, req.Bio); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "profile updated"})
}

func (h *UserHandler) UploadProfileImage(c echo.Context) error {

	userID := int64(c.Get("user_id").(float64))

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "image required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// create uploads folder if not exists
	os.MkdirAll("uploads/profile-images", os.ModePerm)

	filePath := fmt.Sprintf("uploads/profile-images/user-%d.png", userID)

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	imageURL := "/" + filePath

	err = h.userService.UpdateProfileImage(userID, imageURL)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"profile_image": imageURL,
	})
}

func (h *UserHandler) GetActivity(c echo.Context) error {

	userIDParam := c.Param("id")
	userID, _ := strconv.ParseInt(userIDParam, 10, 64)

	requestUserID := int64(c.Get("user_id").(float64))

	// Only allow own activity
	if userID != requestUserID {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "not authorized",
		})
	}

	res, err := h.userService.GetUserActivity(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, res)
}

// get my profile
func (h *UserHandler) GetMyProfile(c echo.Context) error {

	userID := int64(c.Get("user_id").(float64))

	res, err := h.userService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetCommunityStats(c echo.Context) error {
	users, questions, answers, err := h.userService.GetCommunityStats()
	if err != nil {
		return c.JSON(500, echo.Map{"error": "failed"})
	}

	return c.JSON(200, echo.Map{
		"total_users":     users,
		"total_questions": questions,
		"total_answers":   answers,
	})
}

func (h *UserHandler) GetUsers(c echo.Context) error {

	pageParam := c.QueryParam("page")
	if pageParam == "" {
		pageParam = "1"
	}

	page, _ := strconv.Atoi(pageParam)

	limit := 20

	res, err := h.userService.GetUsers(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, res)
}