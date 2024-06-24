package auth

import (
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type BCryptLogin struct {
	cost *int
}

func NewBCryptLogin(cost *int) *BCryptLogin {
	return &BCryptLogin{
		cost: cost,
	}
}

func (b *BCryptLogin) CreatePassword(password *string, hashedPassword *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), *b.cost)
	if err != nil {
		return err
	}
	*hashedPassword = string(hash[:])
	return nil
}

func (b *BCryptLogin) VerifyPassword(input *string, user *store.User) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*input))
}
