package repositories

import (
	"errors"
	"time"

	"github.com/sadia-54/qstack-backend/internal/models/domains"

	"gorm.io/gorm"
)

type EmailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) *EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{db}
}

// Create a new email verification token
func (r *EmailVerificationTokenRepository) CreateToken(token *domains.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

// Find a token by its hash (only if not used and not expired)
func (r *EmailVerificationTokenRepository) FindValidToken(tokenHash string) (*domains.EmailVerificationToken, error) {

	var token domains.EmailVerificationToken

	err := r.db.Where("token_hash = ?", tokenHash).
		First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &token, nil
}

// Mark a token as used
func (r *EmailVerificationTokenRepository) MarkTokenUsed(id int64) error {
	return r.db.Model(&domains.EmailVerificationToken{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).Error
}
