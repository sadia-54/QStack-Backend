package repositories

import (
	"errors"
	"time"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"gorm.io/gorm"
)

type PasswordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db}
}

func (r *PasswordResetTokenRepository) CreateToken(token *domains.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *PasswordResetTokenRepository) FindValidToken(hash string) (*domains.PasswordResetToken, error) {

	var token domains.PasswordResetToken

	err := r.db.Where("token_hash = ?", hash).First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if token.UsedAt != nil {
		return nil, nil
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, nil
	}

	return &token, nil
}

func (r *PasswordResetTokenRepository) MarkUsed(id int64) error {
	return r.db.Model(&domains.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", time.Now()).Error
}