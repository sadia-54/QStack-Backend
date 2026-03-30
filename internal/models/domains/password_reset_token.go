package domains

import "time"

type PasswordResetToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	UserID    int64      `gorm:"not null;index"`
	TokenHash string     `gorm:"type:text;not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt time.Time  `gorm:"not null"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

func NewPasswordResetToken(userID int64, tokenHash string, expires time.Time) *PasswordResetToken {
	return &PasswordResetToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expires,
		CreatedAt: time.Now(),
	}
}