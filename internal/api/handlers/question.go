package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/services"
)

type QuestionHandler struct {
	service *services.QuestionService
}

func NewQuestionHandler(s *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: s}
}

func (h *QuestionHandler) Create(c echo.Context) error {

	var req dtos.CreateQuestion
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	userID := c.Get("user_id").(float64)

	res, err := h.service.Create(int64(userID), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *QuestionHandler) Update(c echo.Context) error {

	idParam := c.Param("id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	var req dtos.UpdateQuestion
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	userID := c.Get("user_id").(float64)

	if err := h.service.Update(int64(userID), id, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

func (h *QuestionHandler) Feed(c echo.Context) error {

	search := c.QueryParam("search")
	tag := c.QueryParam("tag")
	sort := c.QueryParam("sort")

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = o
		}
	}

	questions, err := h.service.GetFeed(search, tag, sort, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) MyFeed(c echo.Context) error {

	userID := int64(c.Get("user_id").(float64))

	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err == nil {
			offset = o
		}
	}

	questions, err := h.service.GetMyFeed(userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed"})
	}

	return c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) GetByID(c echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	question, err := h.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "question not found"})
	}

	return c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) Delete(c echo.Context) error {

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	userID := c.Get("user_id").(float64) // JWT claims are float64
	userIDInt := int64(userID)

	if err := h.service.Delete(userIDInt, id); err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}

func (h *QuestionHandler) Vote(c echo.Context) error {

	idParam := c.Param("id")
	questionID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	userID := int64(c.Get("user_id").(float64))

	var req dtos.VoteQuestion

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := h.service.Vote(userID, questionID, req.Value); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "vote updated"})
}