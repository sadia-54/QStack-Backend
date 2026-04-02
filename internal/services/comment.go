package services

import (
	"errors"
	"time"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

type CommentService struct {
	commentRepo *repositories.CommentRepository
	answerRepo  *repositories.AnswerRepository
}

func NewCommentService(cr *repositories.CommentRepository, ar *repositories.AnswerRepository) *CommentService {
	return &CommentService{
		commentRepo: cr,
		answerRepo:  ar,
	}
}

// CREATE COMMENT
func (s *CommentService) Create(userID, answerID int64, req dtos.CreateComment) (*dtos.CommentResponse, error) {

	// ensure answer exists
	_, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return nil, errors.New("answer not found")
	}

	comment := &domains.Comment{
		UserID:     userID,
		ParentType: 2,
		ParentID:   answerID,
		Body:       req.Body,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	created, err := s.commentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	return mapToCommentResponse(created), nil
}

// GET COMMENTS
func (s *CommentService) GetByAnswer(answerID int64) ([]*dtos.CommentResponse, error) {

	comments, err := s.commentRepo.GetByAnswerID(answerID)
	if err != nil {
		return nil, err
	}

	var res []*dtos.CommentResponse

	for _, c := range comments {
		res = append(res, mapToCommentResponse(&c))
	}

	return res, nil
}

// DELETE COMMENT
func (s *CommentService) Delete(userID, commentID int64) error {

	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("not authorized")
	}

	return s.commentRepo.Delete(commentID)
}

// UPDATE COMMENT
func (s *CommentService) Update(userID, commentID int64, req dtos.UpdateComment) error {

	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	// authorization check
	if comment.UserID != userID {
		return errors.New("not authorized")
	}

	if req.Body != nil {
		comment.Body = *req.Body
	}

	return s.commentRepo.Update(comment)
}

// mapper
func mapToCommentResponse(c *domains.Comment) *dtos.CommentResponse {
	return &dtos.CommentResponse{
		ID:   c.ID,
		Body: c.Body,
		Author: dtos.UserSummary{
			ID:       c.User.ID,
			Username: c.User.Username,
		},
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}