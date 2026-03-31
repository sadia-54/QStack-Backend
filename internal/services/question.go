package services

import (
	"errors"
	"time"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

type QuestionService struct {
	questionRepo *repositories.QuestionRepository
	tagRepo      *repositories.TagRepository
	voteRepo     *repositories.QuestionVoteRepository
}

func NewQuestionService(qr *repositories.QuestionRepository, tr *repositories.TagRepository, vr *repositories.QuestionVoteRepository) *QuestionService {
	return &QuestionService{
		questionRepo: qr,
		tagRepo:      tr,
		voteRepo:     vr,
	}
}

func (s *QuestionService) Create(userID int64, req dtos.CreateQuestion) (*dtos.QuestionResponse, error) {

	question := &domains.Question{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	// Handle tags
	var tags []domains.Tag
	for _, tagName := range req.Tags {

		cleanTag := domains.NewTag(tagName)

		existing, err := s.tagRepo.FindByName(cleanTag.Name)
		if err == nil {
			tags = append(tags, *existing)
			continue
		}

		// If not found → create
		if err := s.tagRepo.Create(cleanTag); err != nil {
			return nil, err
		}

		tags = append(tags, *cleanTag)
	}

	question.Tags = tags

	if err := s.questionRepo.Create(question); err != nil {
		return nil, err
	}

	// reload with relations
	createdQuestion, err := s.questionRepo.FindByID(question.ID)
	if err != nil {
		return nil, err
	}

	return mapToQuestionResponse(createdQuestion), nil
}

func (s *QuestionService) Update(userID, questionID int64, req dtos.UpdateQuestion) error {

	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		return err
	}

	if question.UserID != userID {
		return errors.New("not authorized")
	}

	if req.Title != nil {
		question.Title = *req.Title
	}

	if req.Description != nil {
		question.Description = *req.Description
	}

	question.UpdatedAt = time.Now()

	return s.questionRepo.Update(question)
}

func (s *QuestionService) GetFeed(
	search string,
	tag string,
	sort string,
	limit int,
	offset int,
) ([]*dtos.QuestionResponse, error) {

	questions, err := s.questionRepo.GetFeed(search, tag, sort, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dtos.QuestionResponse
	for _, q := range questions {
		response = append(response, mapToQuestionResponse(&q))
	}

	return response, nil
}

// for my interested feed 
func (s *QuestionService) GetMyFeed(
	userID int64,
	limit int,
	offset int,
) ([]*dtos.QuestionResponse, error) {

	questions, err := s.questionRepo.GetMyFeed(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dtos.QuestionResponse

	for _, q := range questions {
		response = append(response, mapToQuestionResponse(&q))
	}

	return response, nil
}

func (s *QuestionService) GetByID(id int64) (*dtos.QuestionResponse, error) {

	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return mapToQuestionResponse(question), nil
}

func (s *QuestionService) Delete(userID, questionID int64) error {

	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		return err
	}

	if question.UserID != userID {
		return errors.New("not authorized")
	}

	return s.questionRepo.Delete(questionID)
}

func mapToQuestionResponse(q *domains.Question) *dtos.QuestionResponse {

	var tagNames []string
	for _, t := range q.Tags {
		tagNames = append(tagNames, t.Name)
	}

	return &dtos.QuestionResponse{
		ID:          q.ID,
		Title:       q.Title,
		Description: q.Description,
		VoteCount:   q.VoteCount,
		AnswerCount: q.AnswerCount,
		Author: dtos.UserSummary{
			ID:       q.User.ID,
			Username: q.User.Username,
		},
		Tags:      tagNames,
		CreatedAt: q.CreatedAt.Format(time.RFC3339),
		UpdatedAt: q.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *QuestionService) Vote(userID, questionID int64, value int) error {

	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		return err
	}

	// owner cannot vote
	if question.UserID == userID {
		return errors.New("cannot vote own question")
	}

	existing, err := s.voteRepo.Find(questionID, userID)

	// first vote
	if err != nil {

		if !s.voteRepo.IsNotFound(err) {
			return err
		}

		vote := &domains.QuestionVote{
			QuestionID: questionID,
			UserID:     userID,
			Value:      value,
		}

		if err := s.voteRepo.Create(vote); err != nil {
			return err
		}

		return s.questionRepo.UpdateVoteCount(questionID, value)
	}

	// same vote → remove vote
	if existing.Value == value {

		if err := s.voteRepo.Delete(existing.ID); err != nil {
			return err
		}

		return s.questionRepo.UpdateVoteCount(questionID, -value)
	}

	// change vote
	diff := value - existing.Value

	existing.Value = value

	if err := s.voteRepo.Update(existing); err != nil {
		return err
	}

	return s.questionRepo.UpdateVoteCount(questionID, diff)
}

// user owned questions for profile page
func (s *QuestionService) GetMyQuestions(userID int64, limit int, offset int) ([]*dtos.QuestionResponse, error) {

	questions, err := s.questionRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []*dtos.QuestionResponse

	for _, q := range questions {
		response = append(response, mapToQuestionResponse(&q))
	}

	return response, nil
}

func (s *QuestionService) GetPopularTags(limit int) ([]dtos.TagStat, error) {
	return s.questionRepo.GetPopularTags(limit)
}