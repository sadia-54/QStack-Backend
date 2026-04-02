package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/sadia-54/qstack-backend/internal/models/dtos"
	"github.com/sadia-54/qstack-backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// --------------------------------------------
// POST /auth/signup
// --------------------------------------------
func (h *AuthHandler) Signup(c echo.Context) error {
	var body dtos.Signup
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	verifyURL, err := h.authService.Signup(body.Email, body.Username, body.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":     "Signup successful. Please verify your email.",
		"verify_url":  verifyURL, // temporary until RabbitMQ/Mailpit
	})
}

// --------------------------------------------
// POST /auth/login
// --------------------------------------------
func (h *AuthHandler) Login(c echo.Context) error {
	var body dtos.Login
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := h.authService.Login(body.Identifier, body.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// ------------------------------
	// Set access token cookie
	// ------------------------------
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,                    // JS can't access it
		Secure:   false,                   // true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,    // or StrictMode
		MaxAge:   60 * 60 * 24,            // 1 day
	}
	c.SetCookie(accessCookie)

	// ------------------------------
	// Set refresh token cookie
	// ------------------------------
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	}
	c.SetCookie(refreshCookie)

	// ------------------------------
	// Return optional JSON response
	// ------------------------------
	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
	})
}

// logout - clear cookies
func (h *AuthHandler) Logout(c echo.Context) error {
	// Clear access_token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(accessCookie)

	// Clear refresh_token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(refreshCookie)

	return c.JSON(200, echo.Map{
		"message": "logged out",
	})
}

// --------------------------------------------
// GET /auth/verify-email?token=xxx
// --------------------------------------------
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	rawToken := c.QueryParam("token")
	if rawToken == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing token"})
	}

	if err := h.authService.VerifyEmail(rawToken); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Email verified successfully",
	})
}

// change password
func (h *AuthHandler) ChangePassword(c echo.Context) error {

	userID := int64(c.Get("user_id").(float64))

	var body dtos.ChangePassword

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := h.authService.ChangePassword(userID, body.CurrentPassword, body.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "password changed successfully",
	})
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {

	var body dtos.ForgotPasswordRequest

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	h.authService.ForgotPassword(body.Email)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "If the email exists, reset link was sent",
	})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {

	var body dtos.ResetPasswordRequest

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := h.authService.ResetPassword(body.Token, body.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "password reset successful",
	})
}