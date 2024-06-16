package sqlite

import (
	"github.com/google/uuid"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"gorm.io/gorm"
)

type SqliteWebhookConfigStore struct {
	db *gorm.DB
}

func NewSqliteWebhookStore(db *gorm.DB) *SqliteWebhookConfigStore {
	return &SqliteWebhookConfigStore{db: db}
}

func (s *SqliteWebhookConfigStore) CreateWebhookConfig(config *store.WebhookConfig) error {
	return s.db.Create(&config).Error
}

func (s *SqliteWebhookConfigStore) GetWebhookConfig(id uuid.UUID, config *store.WebhookConfig) error {
	err := s.db.Where("id = ?", id).First(&config).Error
	return err
}

func (s *SqliteWebhookConfigStore) GetAllWebhookConfigs(configs *[]store.WebhookConfig) error {
	return s.db.Find(&configs).Error
}

func (s *SqliteWebhookConfigStore) UpdateWebhookConfig(config *store.WebhookConfig) error {
	return s.db.Save(&config).Error
}

func (s *SqliteWebhookConfigStore) DeleteWebhookConfig(id uuid.UUID) error {
	return s.db.Delete(&store.WebhookConfig{}, id).Error
}
