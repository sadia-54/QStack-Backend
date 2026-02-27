package domains

import "time"

type EmailVerificationToken struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"`
	TokenHash string    `gorm:"type:text;not null"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	UsedAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

func NewEmailVerificationToken(userID int64, tokenHash string, expiresAt time.Time) *EmailVerificationToken {
	return &EmailVerificationToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}