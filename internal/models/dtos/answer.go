package dtos

type CreateAnswer struct {
	Description string `json:"description" validate:"required,min=10"`
}

type UpdateAnswer struct {
	Description *string `json:"description,omitempty"`
}

type AnswerResponse struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`

	IsAccepted  bool   `json:"is_accepted"`

	Author      UserSummary `json:"author"`

	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}