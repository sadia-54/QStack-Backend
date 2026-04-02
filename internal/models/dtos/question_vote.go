package dtos

import "time"

type VoteQuestion struct {
	Value int `json:"value" validate:"required,oneof=1 -1"`
}

type QuestionVoteWithDetails struct {
	ID            int64     `json:"id"`
	QuestionID    int64     `json:"question_id"`
	Value         int       `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	QuestionTitle string    `json:"question_title"`
}
