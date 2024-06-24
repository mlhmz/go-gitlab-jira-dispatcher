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

type UserStore interface {
	CreateUser(user *User) error
	GetUser(username string, user *User) error
	GetAllUsers(users *[]User) error
	DeleteUser(username string) error
}
