package domains

import "time"

type Answer struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	QuestionID  int64     `gorm:"not null;index"`
	UserID      int64     `gorm:"not null;index"`

	Description string    `gorm:"type:text;not null"`

	IsAccepted  bool      `gorm:"not null;default:false"`

	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// Relations
	User     *User     `gorm:"foreignKey:UserID"`
	Question *Question `gorm:"foreignKey:QuestionID"`
}

func (Answer) TableName() string {
	return "answers"
}