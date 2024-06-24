package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/auth"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store/sqlite"
	"gorm.io/gorm"
)

func AddAuthToApp(app *fiber.App, db *gorm.DB, token auth.Token) {
	cost := 10
	var login auth.Login = auth.NewBCryptLogin(&cost)

	var userStore store.UserStore = sqlite.NewSqliteUserStore(db)

	app.Post("/login", func(c *fiber.Ctx) error {
		var payload auth.LoginPayload

		if err := c.BodyParser(&payload); err != nil {
			payloadErr := fmt.Errorf("failed to parse payload: %s", err)
			log.Error(payloadErr)
			return c.Status(500).SendString(payloadErr.Error())
		}

		var user store.User
		authFailMsg := fmt.Sprintf("failed authorization with the user '%s'", payload.Username)
		if err := userStore.GetUser(payload.Username, &user); err != nil {
			log.Error(err)
			return c.Status(401).SendString(authFailMsg)
		}

		if err := login.VerifyPassword(&payload.Password, &user); err != nil {
			log.Error(err)
			return c.Status(401).SendString(authFailMsg)
		}

		var tokenStr string
		if err := token.CreateToken(&user.Username, &tokenStr); err != nil {
			log.Error(err)
			c.Response().Header.Add("WWW-Authenticate", err.Error())
			return c.SendStatus(401)
		}

		c.Cookie(&fiber.Cookie{
			Name:  "token",
			Value: tokenStr,
		})

		return c.SendStatus(200)
	})
}

func AuthMiddleware(token auth.Token) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if err := token.VerifyToken(&tokenStr); err != nil {
			c.Response().Header.Add("WWW-Authenticate", err.Error())
			c.SendStatus(401)
		}
		return c.Next()
	}
}
