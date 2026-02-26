package repositories

import (
	"github.com/sadia-54/qstack-backend/internal/models/domains"
	gormmodels "github.com/sadia-54/qstack-backend/internal/models/gorm"

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
	gormToken := mapDomainToGormEmailToken(token)
	return r.db.Create(&gormToken).Error
}

// Find a token by its hash (only if not used and not expired)
func (r *EmailVerificationTokenRepository) FindValidToken(tokenHash string) (*domains.EmailVerificationToken, error) {
	var gormToken gormmodels.EmailVerificationToken

	err := r.db.Where(
		"token_hash = ? AND used_at IS NULL AND expires_at > NOW()",
		tokenHash,
	).First(&gormToken).Error

	if err != nil {
		return nil, err
	}

	return mapGormToDomainEmailToken(&gormToken), nil
}

// Mark a token as used
func (r *EmailVerificationTokenRepository) MarkTokenUsed(id int64) error {
	return r.db.Model(&gormmodels.EmailVerificationToken{}).
		Where("id = ?", id).
		Update("used_at", gorm.Expr("NOW()")).Error
}

// Mapping: Domain -> GORM model
func mapDomainToGormEmailToken(t *domains.EmailVerificationToken) gormmodels.EmailVerificationToken {
	return gormmodels.EmailVerificationToken{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
		UsedAt:    t.UsedAt,
		CreatedAt: t.CreatedAt,
	}
}

// Mapping: GORM model -> Domain model
func mapGormToDomainEmailToken(t *gormmodels.EmailVerificationToken) *domains.EmailVerificationToken {
	return &domains.EmailVerificationToken{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
		UsedAt:    t.UsedAt,
		CreatedAt: t.CreatedAt,
	}
}