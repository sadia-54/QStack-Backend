package handlers

import (
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