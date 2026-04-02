package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/sadia-54/qstack-backend/internal/models/domains"
	"github.com/sadia-54/qstack-backend/internal/queue"
	"github.com/sadia-54/qstack-backend/internal/repositories"
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	tokenRepo  *repositories.EmailVerificationTokenRepository
	resetRepo  *repositories.PasswordResetTokenRepository
	jwtSecret  string
	appBaseURL string
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	tokenRepo *repositories.EmailVerificationTokenRepository,
	resetRepo *repositories.PasswordResetTokenRepository,
	jwtSecret string,
	appBaseURL string, // e.g., http://localhost:3000
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		resetRepo:  resetRepo,
		jwtSecret:  jwtSecret,
		appBaseURL: appBaseURL,
	}
}

// SIGNUP
func (s *AuthService) Signup(email, username, password string) (string, error) {
	// 1. Check if user exists already
	existing, err := s.userRepo.FindByEmailOrUsername(email)
	if err != nil {
		return "", errors.New("database error")
	}
	if existing != nil {
		return "", errors.New("email already registered")
	}

	existing, err = s.userRepo.FindByEmailOrUsername(username)
	if err != nil {
		return "", errors.New("database error")
	}
	if existing != nil {
		return "", errors.New("username already taken")
	}

	// 2. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 3. Create domain user
	user := domains.NewUser(email, username, string(hashed))

	// 4. Save user
	if err := s.userRepo.CreateUser(user); err != nil {
		return "", err
	}

	// 5. Create verification token
	rawToken, tokenHash, expiresAt, err := generateEmailVerificationToken()
	if err != nil {
		return "", err
	}

	token := domains.NewEmailVerificationToken(user.ID, tokenHash, expiresAt)

	// save token
	if err := s.tokenRepo.CreateToken(token); err != nil {
		return "", err
	}

	// 6. Return temporary verification link 
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.appBaseURL, url.QueryEscape(rawToken))

	queue.PublishEmailVerification(email, rawToken)

	return verifyURL, nil
}

// LOGIN
func (s *AuthService) Login(identifier, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmailOrUsername(identifier)
	if err != nil {
		return "", "", errors.New("database error")
	}

	if user == nil {
		return "", "", errors.New("invalid credentials")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Require email verification
	if !user.EmailVerified {
		return "", "", errors.New("email not verified")
	}

	// Generate JWT tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// EMAIL VERIFICATION
func (s *AuthService) VerifyEmail(rawToken string) error {
	// Hash raw token
	tokenHash := hashToken(rawToken)

	// Find valid token
	token, err := s.tokenRepo.FindValidToken(tokenHash)
	if err != nil {
		return errors.New("Database Error")
	}

	if token == nil {
		return errors.New("Invalid or expired token")
	}

	// // check if already used
	// if token.UsedAt != nil {
	// 	return errors.New("Token already used")
	// }

	// check if expired
	if time.Now().After(token.ExpiresAt) {
		return errors.New("Token expired")
	}

	// Mark token as used
	if err := s.tokenRepo.MarkTokenUsed(token.ID); err != nil {
		return err
	}

	// Update user email_verified = true
	user, err := s.userRepo.GetUserByID(token.UserID)
	if err != nil {
		return err
	}

	user.EmailVerified = true
	user.UpdatedAt = time.Now()

	return s.userRepo.UpdateUser(user)
}

// JWT GENERATION
func (s *AuthService) generateAccessToken(user *domains.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) generateRefreshToken(user *domains.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// TOKEN GENERATION HELPERS
func generateEmailVerificationToken() (rawToken string, hashed string, expires time.Time, err error) {
	// Create 32-byte secure random token
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return "", "", time.Time{}, err
	}

	raw := base64.RawURLEncoding.EncodeToString(b)
	hashed = hashToken(raw)
	expires = time.Now().Add(24 * time.Hour)

	return raw, hashed, expires, nil
}

func hashToken(token string) string {
	// Simple hash using bcrypt or SHA256 (bcrypt optional because it's slow)
	// We choose SHA256 for speed
	h := sha256Sum(token)
	return h
}

func sha256Sum(s string) string {
    b := sha256.Sum256([]byte(s))
    return base64.RawURLEncoding.EncodeToString(b[:])
}

// change password
func (s *AuthService) ChangePassword(userID int64, currentPassword, newPassword string) error {

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// verify current password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(currentPassword),
	); err != nil {
		return errors.New("current password incorrect")
	}

	// prevent reusing same password
	if currentPassword == newPassword {
		return errors.New("new password must be different")
	}

	// hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)
	user.UpdatedAt = time.Now()

	return s.userRepo.UpdateUser(user)
}

func (s *AuthService) ForgotPassword(email string) error {

	user, err := s.userRepo.FindByEmailOrUsername(email)
	if err != nil || user == nil {
		return nil // do not reveal user existence
	}

	rawToken, tokenHash, expires, err := generateEmailVerificationToken()
	if err != nil {
		return err
	}

	token := domains.NewPasswordResetToken(user.ID, tokenHash, expires)

	if err := s.resetRepo.CreateToken(token); err != nil {
		return err
	}

	return queue.PublishPasswordReset(user.Email, rawToken)
}

func (s *AuthService) ResetPassword(rawToken, newPassword string) error {

	hash := hashToken(rawToken)

	token, err := s.resetRepo.FindValidToken(hash)
	if err != nil || token == nil {
		return errors.New("invalid or expired token")
	}

	user, err := s.userRepo.GetUserByID(token.UserID)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashed)
	user.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return s.resetRepo.MarkUsed(token.ID)
}