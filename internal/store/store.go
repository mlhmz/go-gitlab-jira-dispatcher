package store

import "github.com/google/uuid"

type WebhookStore interface {
	CreateWebhook(webhook Webhook) error
	GetWebhook(id uuid.UUID) (Webhook, error)
	UpdateWebhook(webhook Webhook) error
	DeleteWebhook(id uuid.UUID) error
}
