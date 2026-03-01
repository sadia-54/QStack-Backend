package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *domains.Question) error {
	return r.db.Create(question).Error
}

func (r *QuestionRepository) Update(question *domains.Question) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(question).Error
}

func (r *QuestionRepository) FindByID(id int64) (*domains.Question, error) {
	var question domains.Question
	err := r.db.Preload("User").
		Preload("Tags").
		First(&question, id).Error

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) GetFeed(limit, offset int) ([]domains.Question, error) {
	var questions []domains.Question
	err := r.db.Preload("User").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&questions).Error

	return questions, err
}

func (r *QuestionRepository) Delete(id int64) error {
	return r.db.Delete(&domains.Question{}, id).Error
}