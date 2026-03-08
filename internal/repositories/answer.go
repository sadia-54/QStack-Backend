package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(answer *domains.Answer) error {
	return r.db.Create(answer).Error
}

func (r *AnswerRepository) Update(answer *domains.Answer) error {
	return r.db.Save(answer).Error
}

func (r *AnswerRepository) FindByID(id int64) (*domains.Answer, error) {

	var answer domains.Answer

	err := r.db.
		Preload("User").
		First(&answer, id).Error

	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *AnswerRepository) GetByQuestionID(questionID int64) ([]domains.Answer, error) {

	var answers []domains.Answer

	err := r.db.
		Preload("User").
		Where("question_id = ?", questionID).
		Order("is_accepted DESC").
		Order("created_at ASC").
		Find(&answers).Error

	return answers, err
}

func (r *AnswerRepository) Delete(id int64) error {
	return r.db.Delete(&domains.Answer{}, id).Error
}

func (r *AnswerRepository) AcceptAnswer(answerID int64) error {

	return r.db.Model(&domains.Answer{}).
		Where("id = ?", answerID).
		Update("is_accepted", true).Error
}