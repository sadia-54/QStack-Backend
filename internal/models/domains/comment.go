package domains

import "time"

type Comment struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null;index"`
	ParentType int16     `gorm:"not null"` // 2 = answer
	ParentID   int64     `gorm:"not null;index"`

	Body       string    `gorm:"type:varchar(1000);not null"`

	CreatedAt  time.Time `gorm:"type:timestamptz;default:now()"`

	// Relations
	User *User `gorm:"foreignKey:UserID"`
}

func (Comment) TableName() string {
	return "comments"
}