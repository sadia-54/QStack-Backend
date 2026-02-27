package repositories

import (
	"errors"

	"github.com/sadia-54/qstack-backend/internal/models/domains"

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
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmailOrUsername (identifier string) (*domains.User, error) {
	var user domains.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *UserRepository) UpdateUser (user *domains.User) error {

	return r.db.Save(user).Error
}

func (r *UserRepository) GetUserByID(id int64) (*domains.User, error){
	var user domains.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}