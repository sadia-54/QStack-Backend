package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Find by name
func (r *TagRepository) FindByName(name string) (*domains.Tag, error) {
	var tag domains.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create tag
func (r *TagRepository) Create(tag *domains.Tag) error {
	return r.db.Create(tag).Error
}