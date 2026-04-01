package dtos

// ----------- AUTH REQUESTS -----------

type Signup struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type Login struct {
	Identifier string `json:"identifier" validate:"required"` // can be email OR username
	Password   string `json:"password" validate:"required"`
}

// ----------- AUTH RESPONSES -----------

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// ----------- USER PROFILE DTO -----------

type UserDTO struct {
	ID                        int64   `json:"id"`
	Email                     string  `json:"email"`
	Username                  string  `json:"username"`
	Bio                       *string `json:"bio,omitempty"`
	EmailVerified             bool    `json:"email_verified"`
	EmailNotificationsEnabled bool    `json:"email_notifications_enabled"`
}

type Profile struct {
	ID             int64    `json:"id"`
	Username       string   `json:"username"`
	Bio            string   `json:"bio"`
	TotalQuestions int64    `json:"total_questions"`
	TotalAnswers   int64    `json:"total_answers"`
	TotalVotes     int64    `json:"total_votes"`
	PreferredTags  []string `json:"preferred_tags"`
	CreatedAt      string   `json:"created_at"`
}

// user activity 
type ActivityItem struct {
	Type      string `json:"type"` // question, answer, vote, edit, accept
	Title     string `json:"title,omitempty"`
	TargetID  int64  `json:"target_id,omitempty"`
	Value     int    `json:"value,omitempty"`
	CreatedAt string `json:"created_at"`
}

// password change
type ChangePassword struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// for public user listing 
type UserSummaryPublic struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	TotalQuestions int64  `json:"total_questions"`
	TotalAnswers   int64  `json:"total_answers"`
	TotalVotes     int64  `json:"total_votes"`
	CreatedAt      string `json:"created_at"`
}