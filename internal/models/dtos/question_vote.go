package dtos

type VoteQuestion struct {
	Value int `json:"value" validate:"required,oneof=1 -1"`
}