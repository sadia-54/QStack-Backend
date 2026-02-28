package domains

import (
	"strings"
	"time"
)

type Tag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`

	Questions []Question `gorm:"many2many:question_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}

func NewTag(name string) *Tag {
	return &Tag{
		Name: strings.ToLower(strings.TrimSpace(name)),
	}
}