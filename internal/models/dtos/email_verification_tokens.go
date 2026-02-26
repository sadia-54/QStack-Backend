package dtos

type VerifyEmailRequest struct {
	Token string `query:"token" validate:"required"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
}