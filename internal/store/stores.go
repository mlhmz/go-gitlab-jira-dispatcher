package store

import (
	"github.com/google/uuid"
)

type WebhookConfigStore interface {
	CreateWebhookConfig(webhook *WebhookConfig) error
	GetWebhookConfig(id uuid.UUID, webhook *WebhookConfig) error
	GetAllWebhookConfigs(webhooks *[]WebhookConfig) error
	UpdateWebhookConfig(webhook *WebhookConfig) error
	DeleteWebhookConfig(id uuid.UUID) error
}
