package dtos

type CreateComment struct {
	Body string `json:"body" validate:"required,min=2,max=1000"`
}

type CommentResponse struct {
	ID        int64       `json:"id"`
	Body      string      `json:"body"`
	Author    UserSummary `json:"author"`
	CreatedAt string      `json:"created_at"`
}

type UpdateComment struct {
	Body *string `json:"body,omitempty"`
}