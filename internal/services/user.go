package services

import (

	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID int64) (*dtos.Profile, error) {

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	q, a, v, err := s.userRepo.GetProfileStats(userID)
	if err != nil {
		return nil, err
	}

	var bio string
	if user.Bio != nil {
		bio = *user.Bio
	}

	return &dtos.Profile{
		ID:             user.ID,
		Username:       user.Username,
		Bio:            bio,
		TotalQuestions: q,
		TotalAnswers:   a,
		TotalVotes:     v,
	}, nil
}

func (s *UserService) UpdateProfile(userID int64, bio string) error {
	return s.userRepo.UpdateBio(userID, bio)
}