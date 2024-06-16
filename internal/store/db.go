package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&WebhookConfig{}); err != nil {
		panic(err)
	}

	return db
}
