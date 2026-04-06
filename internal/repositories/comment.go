package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *domains.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByAnswerID(answerID int64) ([]domains.Comment, error) {
	var comments []domains.Comment

	err := r.db.
		Preload("User").
		Where("parent_type = ? AND parent_id = ?", 2, answerID).
		Order("created_at ASC").
		Find(&comments).Error

	return comments, err
}

func (r *CommentRepository) FindByID(id int64) (*domains.Comment, error) {
	var c domains.Comment

	err := r.db.Preload("User").First(&c, id).Error
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CommentRepository) Delete(id int64) error {
	return r.db.Delete(&domains.Comment{}, id).Error
}

func (r *CommentRepository) Update(comment *domains.Comment) error {
	return r.db.Save(comment).Error
}
