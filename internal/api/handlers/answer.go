package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/services"
)

type AnswerHandler struct {
	service *services.AnswerService
}

func NewAnswerHandler(s *services.AnswerService) *AnswerHandler {
	return &AnswerHandler{service: s}
}

func (h *AnswerHandler) Create(c echo.Context) error {

	questionID, _ := strconv.ParseInt(c.Param("question_id"), 10, 64)

	var req dtos.CreateAnswer
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	userID := c.Get("user_id").(float64)

	res, err := h.service.Create(int64(userID), questionID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *AnswerHandler) GetByQuestion(c echo.Context) error {

	questionID, _ := strconv.ParseInt(c.Param("question_id"), 10, 64)

	res, err := h.service.GetByQuestion(questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AnswerHandler) Update(c echo.Context) error {

	answerID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req dtos.UpdateAnswer
	c.Bind(&req)

	userID := c.Get("user_id").(float64)

	err := h.service.Update(int64(userID), answerID, req)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

func (h *AnswerHandler) Delete(c echo.Context) error {

	answerID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	userID := c.Get("user_id").(float64)

	err := h.service.Delete(int64(userID), answerID)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}

func (h *AnswerHandler) Accept(c echo.Context) error {

	answerID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	userID := c.Get("user_id").(float64)

	err := h.service.AcceptAnswer(int64(userID), answerID)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "accepted"})
}