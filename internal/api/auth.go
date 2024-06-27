package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/auth"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

type AuthApi struct {
	login     auth.Login
	userStore store.UserStore
	token     auth.Token
	app       *fiber.App
}

func NewAuthApi(login auth.Login, userStore store.UserStore, token auth.Token, app *fiber.App) *AuthApi {
	return &AuthApi{
		login:     login,
		userStore: userStore,
		token:     token,
		app:       app,
	}
}

func (authApi *AuthApi) AddAuthToApp(app *fiber.App) {
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}, "layouts/main")
	})

	app.Post("/authenticate", func(c *fiber.Ctx) error {
		if _, err := authApi.authenticateAndIssueToken(c); err != nil {
			return c.Render("login", fiber.Map{
				"ErrorMsg": err.Error(),
			})
		}
		return c.Redirect("/")
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		c.ClearCookie("token")
		return c.Redirect("/login")
	})

	app.Post("/api/login", func(c *fiber.Ctx) error {
		if status, err := authApi.authenticateAndIssueToken(c); err != nil {
			return c.Status(status).SendString(err.Error())
		}

		return c.SendStatus(200)
	})
}

func (authApi *AuthApi) GetAuthRedirectMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if err := authApi.token.VerifyToken(&tokenStr); err != nil {
			c.Response().Header.Add("WWW-Authenticate", err.Error())
			return c.Redirect("/login")
		}
		return c.Next()
	}
}

func (authApi *AuthApi) GetAuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if err := authApi.token.VerifyToken(&tokenStr); err != nil {
			c.Response().Header.Add("WWW-Authenticate", err.Error())
			return c.SendStatus(401)
		}
		return c.Next()
	}
}

func (authApi *AuthApi) authenticateAndIssueToken(c *fiber.Ctx) (int, error) {
	var payload auth.LoginPayload
	var accessToken string

	if err := c.BodyParser(&payload); err != nil {
		payloadErr := fmt.Errorf("failed to parse payload: %s", err)
		log.Error(payloadErr)
		return 500, payloadErr
	}

	var user store.User
	authFailMsg := fmt.Errorf("failed authorization with the user '%s'", payload.Username)
	if err := authApi.userStore.GetUser(payload.Username, &user); err != nil {
		log.Error(err)
		return 401, authFailMsg
	}

	if err := authApi.login.VerifyPassword(&payload.Password, &user); err != nil {
		log.Error(err)
		return 401, authFailMsg
	}

	if err := authApi.token.CreateToken(&user.Username, &accessToken); err != nil {
		log.Error(err)
		c.Response().Header.Add("WWW-Authenticate", err.Error())
		return 401, err
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: accessToken,
	})

	return 200, nil
}
