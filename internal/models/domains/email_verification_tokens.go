package domains

import "time"

type EmailVerificationToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

func NewEmailVerificationToken(userID int64, tokenHash string, expiresAt time.Time) *EmailVerificationToken {
	return &EmailVerificationToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}