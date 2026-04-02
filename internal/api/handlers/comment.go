package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/services"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(s *services.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

// CREATE
func (h *CommentHandler) Create(c echo.Context) error {

	answerID, _ := strconv.ParseInt(c.Param("answer_id"), 10, 64)

	var req dtos.CreateComment
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}

	userID := c.Get("user_id").(float64)

	res, err := h.service.Create(int64(userID), answerID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

// GET
func (h *CommentHandler) GetByAnswer(c echo.Context) error {

	answerID, _ := strconv.ParseInt(c.Param("answer_id"), 10, 64)

	res, err := h.service.GetByAnswer(answerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, res)
}

// DELETE
func (h *CommentHandler) Delete(c echo.Context) error {

	commentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	userID := c.Get("user_id").(float64)

	err := h.service.Delete(int64(userID), commentID)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}

// UPDATE
func (h *CommentHandler) Update(c echo.Context) error {

	commentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req dtos.UpdateComment
	c.Bind(&req)

	userID := c.Get("user_id").(float64)

	err := h.service.Update(int64(userID), commentID, req)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}