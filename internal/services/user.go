package services

import (
	"time"
	"sort"

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

	var profileImage string
	if user.ProfileImage != nil {
		profileImage = *user.ProfileImage
	}

	return &dtos.Profile{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Bio:            bio,
		ProfileImage:   profileImage,
		TotalQuestions: q,
		TotalAnswers:   a,
		TotalVotes:     v,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) UpdateProfile(userID int64, bio string) error {
	return s.userRepo.UpdateBio(userID, bio)
}

func (s *UserService) UpdateProfileImage(userID int64, imagePath string) error {
	return s.userRepo.UpdateProfileImage(userID, imagePath)
}

// user activity 
func (s *UserService) GetUserActivity(userID int64) ([]dtos.ActivityItem, error) {

	var activities []dtos.ActivityItem

	// QUESTIONS
	questions, _ := s.userRepo.GetUserQuestions(userID, 5)
	for _, q := range questions {
		activities = append(activities, dtos.ActivityItem{
			Type:      "question",
			Title:     q.Title,
			TargetID:  q.ID,
			CreatedAt: q.CreatedAt.Format(time.RFC3339),
		})
	}

	// ANSWERS
	answers, _ := s.userRepo.GetUserAnswers(userID, 5)
	for _, a := range answers {
		activities = append(activities, dtos.ActivityItem{
			Type:      "answer",
			TargetID:  a.QuestionID,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		})
	}

	// VOTES
	votes, _ := s.userRepo.GetUserVotes(userID, 5)
	for _, v := range votes {
		activities = append(activities, dtos.ActivityItem{
			Type:      "vote",
			TargetID:  v.QuestionID,
			Value:     v.Value,
			CreatedAt: v.CreatedAt.Format(time.RFC3339),
		})
	}

	// ACCEPTED
	accepted, _ := s.userRepo.GetAcceptedAnswers(userID, 5)
	for _, a := range accepted {
		activities = append(activities, dtos.ActivityItem{
			Type:      "accept",
			TargetID:  a.QuestionID,
			CreatedAt: a.UpdatedAt.Format(time.RFC3339),
		})
	}

	// EDITED QUESTIONS
	editedQ, _ := s.userRepo.GetEditedQuestions(userID, 5)
	for _, q := range editedQ {
		activities = append(activities, dtos.ActivityItem{
			Type:      "edit",
			TargetID:  q.ID,
			CreatedAt: q.UpdatedAt.Format(time.RFC3339),
		})
	}

	// EDITED ANSWERS
	editedA, _ := s.userRepo.GetEditedAnswers(userID, 5)
	for _, a := range editedA {
		activities = append(activities, dtos.ActivityItem{
			Type:      "edit",
			TargetID:  a.QuestionID,
			CreatedAt: a.UpdatedAt.Format(time.RFC3339),
		})
	}

	// SORT 
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].CreatedAt > activities[j].CreatedAt
	})

	return activities, nil
}

func (s *UserService) GetCommunityStats() (int64, int64, int64, error) {
	return s.userRepo.GetCommunityStats()
}

func (s *UserService) GetUsers(page int, limit int) ([]dtos.UserSummaryPublic, error) {

	offset := (page - 1) * limit

	users, err := s.userRepo.GetUsers(limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil
}