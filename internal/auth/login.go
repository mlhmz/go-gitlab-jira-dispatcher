package auth

import "github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"

type Login interface {
	CreatePassword(password *string, hashedPassword *string) error
	VerifyPassword(input *string, user *store.User) error
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
