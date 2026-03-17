package dtos

type CreateQuestion struct {
	Title       string   `json:"title" validate:"required,min=10,max=200"`
	Description string   `json:"description" validate:"required,min=20"`
	Tags        []string `json:"tags" validate:"required,min=1,dive,required"`
}

type UpdateQuestion struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

type QuestionResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`

	VoteCount   int      `json:"vote_count"`
	AnswerCount int      `json:"answer_count"`

	Author      UserSummary `json:"author"`
	Tags        []string    `json:"tags"`

	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type UserSummary struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}