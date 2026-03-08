package services

import (
	"errors"
	"time"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

type AnswerService struct {
	answerRepo   *repositories.AnswerRepository
	questionRepo *repositories.QuestionRepository
}

func NewAnswerService(ar *repositories.AnswerRepository, qr *repositories.QuestionRepository) *AnswerService {
	return &AnswerService{
		answerRepo:   ar,
		questionRepo: qr,
	}
}

func (s *AnswerService) Create(userID, questionID int64, req dtos.CreateAnswer) (*dtos.AnswerResponse, error) {

	answer := &domains.Answer{
		UserID:      userID,
		QuestionID:  questionID,
		Description: req.Description,
	}

	if err := s.answerRepo.Create(answer); err != nil {
		return nil, err
	}

	// increment answer count
	question, err := s.questionRepo.FindByID(questionID)
	if err == nil {
		question.AnswerCount++
		s.questionRepo.Update(question)
	}

	created, err := s.answerRepo.FindByID(answer.ID)
	if err != nil {
		return nil, err
	}

	return mapToAnswerResponse(created), nil
}

func (s *AnswerService) Update(userID, answerID int64, req dtos.UpdateAnswer) error {

	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return err
	}

	if answer.UserID != userID {
		return errors.New("not authorized")
	}

	if req.Description != nil {
		answer.Description = *req.Description
	}

	answer.UpdatedAt = time.Now()

	return s.answerRepo.Update(answer)
}

func (s *AnswerService) GetByQuestion(questionID int64) ([]*dtos.AnswerResponse, error) {

	answers, err := s.answerRepo.GetByQuestionID(questionID)
	if err != nil {
		return nil, err
	}

	var response []*dtos.AnswerResponse

	for _, a := range answers {
		response = append(response, mapToAnswerResponse(&a))
	}

	return response, nil
}

func (s *AnswerService) Delete(userID, answerID int64) error {

	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return err
	}

	if answer.UserID != userID {
		return errors.New("not authorized")
	}

	err = s.answerRepo.Delete(answerID)
	if err != nil {
		return err
	}

	// decrease answer count
	question, err := s.questionRepo.FindByID(answer.QuestionID)
	if err == nil {
		question.AnswerCount--
		s.questionRepo.Update(question)
	}

	return nil
}

func (s *AnswerService) AcceptAnswer(userID, answerID int64) error {

	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return err
	}

	question, err := s.questionRepo.FindByID(answer.QuestionID)
	if err != nil {
		return err
	}

	if question.UserID != userID {
		return errors.New("only question owner can accept answer")
	}

	return s.answerRepo.AcceptAnswer(answerID)
}

func mapToAnswerResponse(a *domains.Answer) *dtos.AnswerResponse {

	return &dtos.AnswerResponse{
		ID:          a.ID,
		Description: a.Description,
		IsAccepted:  a.IsAccepted,
		Author: dtos.UserSummary{
			ID:       a.User.ID,
			Username: a.User.Username,
		},
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}
}