package store

import (
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookConfigSubmission struct {
	Title           string `json:"title"`
	ReadyForReview  string `json:"ready_for_review"`
	InReview        string `json:"in_review"`
	ReviewOK        string `json:"review_ok"`
	ReviewNotOK     string `json:"review_not_ok"`
	DevelopmentDone string `json:"development_done"`
	Projects        string `json:"projects"`
}

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

func MapSubmission(submission *WebhookConfigSubmission) *WebhookConfig {
	return &WebhookConfig{
		Title:           submission.Title,
		Projects:        submission.Projects,
		ReadyForReview:  parseInt(submission.ReadyForReview),
		InReview:        parseInt(submission.InReview),
		ReviewOK:        parseInt(submission.ReviewOK),
		ReviewNotOK:     parseInt(submission.ReviewNotOK),
		DevelopmentDone: parseInt(submission.DevelopmentDone),
	}
}

func parseInt(value string) int {
	result, err := strconv.Atoi(value)

	if err == nil {
		return result
	} else {
		log.Errorf("there was an error while parsing the value %s - defaulting to 0", value)
		return 0
	}
}
