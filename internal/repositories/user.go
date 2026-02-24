package repositories

import (
	"errors"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	gormmodels "github.com/sadia-54/qstack-backend/internal/models/gorm"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// repository methods
func (r *UserRepository) CreateUser(user *domains.User) error {
	gormUser := mapDomainToGorm(user)
	err := r.db.Create(&gormUser).Error
	if err != nil {
		return err
	}

	// Copy DB-generated fields back to domain model
	user.ID = gormUser.ID
	user.CreatedAt = gormUser.CreatedAt
	user.UpdatedAt = gormUser.UpdatedAt

return nil
}

func (r *UserRepository) FindByEmailOrUsername (identifier string) (*domains.User, error) {
	var gormUser gormmodels.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&gormUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mapGormToDomain(&gormUser), nil

}

func (r *UserRepository) UpdateUser (user *domains.User) error {
	gormUser := mapDomainToGorm(user)
	return r.db.Save(&gormUser).Error
}

func (r *UserRepository) GetUserByID(id int64) (*domains.User, error){
	var gormUser gormmodels.User
	err := r.db.First(&gormUser, id).Error
	if err != nil {
		return nil, err
	}

	return mapGormToDomain(&gormUser), nil
}


// Helper functions to map between domain and gorm models
func mapDomainToGorm(user *domains.User) gormmodels.User {
	return gormmodels.User{
		ID:                        user.ID,
		Email:                     user.Email,
		PasswordHash:              user.PasswordHash,
		Username:                  user.Username,
		Bio:                       user.Bio,
		EmailVerified:             user.EmailVerified,
		EmailNotificationsEnabled: user.EmailNotificationsEnabled,
		CreatedAt:                 user.CreatedAt,
		UpdatedAt:                 user.UpdatedAt,
	}
}

func mapGormToDomain(user *gormmodels.User) *domains.User {
	return &domains.User{
		ID:                        user.ID,
		Email:                     user.Email,
		PasswordHash:              user.PasswordHash,
		Username:                  user.Username,
		Bio:                       user.Bio,
		EmailVerified:             user.EmailVerified,
		EmailNotificationsEnabled: user.EmailNotificationsEnabled,
		CreatedAt:                 user.CreatedAt,
		UpdatedAt:                 user.UpdatedAt,
	}
}