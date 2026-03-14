package domains

import "time"

type QuestionVote struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	QuestionID int64     `gorm:"not null;index"`
	UserID     int64     `gorm:"not null;index"`
	Value      int       `gorm:"not null"` // 1 = upvote, -1 = downvote

	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (QuestionVote) TableName() string {
	return "question_votes"
}