package store

import "github.com/google/uuid"

type Webhook struct {
	ID          uuid.UUID   `json:"id"`
	Transitions Transitions `json:"transitions"`
}

type Transitions struct {
	ReadyForReview  int `json:"ready_for_review"`
	InReview        int `json:"in_review"`
	ReviewOK        int `json:"review_ok"`
	ReviewNotOK     int `json:"review_not_ok"`
	DevelopmentDone int `json:"development_done"`
}
