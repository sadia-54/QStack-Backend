package repositories

import (
	"errors"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type QuestionVoteRepository struct {
	db *gorm.DB
}

func NewQuestionVoteRepository(db *gorm.DB) *QuestionVoteRepository {
	return &QuestionVoteRepository{db: db}
}

func (r *QuestionVoteRepository) Find(questionID, userID int64) (*domains.QuestionVote, error) {

	var vote domains.QuestionVote

	err := r.db.
		Where("question_id = ? AND user_id = ?", questionID, userID).
		First(&vote).Error

	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func (r *QuestionVoteRepository) Create(vote *domains.QuestionVote) error {
	return r.db.Create(vote).Error
}

func (r *QuestionVoteRepository) Update(vote *domains.QuestionVote) error {
	return r.db.Save(vote).Error
}

func (r *QuestionVoteRepository) Delete(id int64) error {
	return r.db.Delete(&domains.QuestionVote{}, id).Error
}

func (r *QuestionVoteRepository) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}