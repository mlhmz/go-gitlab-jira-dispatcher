package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookConfig struct {
	gorm.Model      `json:"-"`
	ID              string `json:"id"`
	Title           string `json:"title"`
	ReadyForReview  int    `json:"ready_for_review"`
	InReview        int    `json:"in_review"`
	ReviewOK        int    `json:"review_ok"`
	ReviewNotOK     int    `json:"review_not_ok"`
	DevelopmentDone int    `json:"development_done"`
	Projects        string `json:"projects"`
}

func (config *WebhookConfig) BeforeCreate(tx *gorm.DB) error {
	config.ID = uuid.NewString()
	return nil
}
