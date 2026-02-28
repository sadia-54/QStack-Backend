package domains

import "time"

type Question struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	UserID      int64     `gorm:"not null;index"`
	Title       string    `gorm:"type:varchar(200);not null"`
	Description string    `gorm:"type:text;not null"`

	VoteCount   int       `gorm:"not null;default:0"`
	AnswerCount int       `gorm:"not null;default:0"`

	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`

	// Relations
	User  *User   `gorm:"foreignKey:UserID"`
	Tags  []Tag  `gorm:"many2many:question_tags;"`
}

func (Question) TableName() string {
	return "questions"
}