package sqlite

import (
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"gorm.io/gorm"
)

type SqliteUserStore struct {
	db *gorm.DB
}

func NewSqliteUserStore(db *gorm.DB) *SqliteUserStore {
	return &SqliteUserStore{
		db: db,
	}
}

func (userStore *SqliteUserStore) CreateUser(user *store.User) error {
	return userStore.db.Create(&user).Error
}

func (userStore *SqliteUserStore) GetUser(username string, user *store.User) error {
	return userStore.db.Where("username = ?", username).First(&user).Error
}

func (userStore *SqliteUserStore) GetAllUsers(users *[]store.User) error {
	return userStore.db.Find(&users).Error
}

func (userStore *SqliteUserStore) DeleteUser(username string) error {
	var user store.User
	if err := userStore.GetUser(username, &user); err != nil {
		return err
	}
	return userStore.db.Delete(&store.User{}).Error
}
