package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/api"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/auth"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/jira"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/jira/jirav2"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store/sqlite"
)

func main() {
	templateEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: templateEngine,
	})

	var jiraUrl string
	var jiraApiToken string
	var tokenSignature string

	loadEnvironment(&jiraUrl, &jiraApiToken, &tokenSignature)

	db := store.NewDatabase()
	var webhookStore store.WebhookConfigStore = sqlite.NewSqliteWebhookStore(db)
	var userStore store.UserStore = sqlite.NewSqliteUserStore(db)

	publisher := gitlab.NewPublisher()
	publisher.Register(jira.NewJiraListener(jirav2.NewRestClient(jiraUrl, jiraApiToken)))

	signature := []byte(tokenSignature)
	cost := 10

	var token auth.Token = auth.NewJwtToken(&signature, time.Hour*1)
	var login auth.Login = auth.NewBCryptLogin(&cost)
	authApi := api.NewAuthApi(login, userStore, token, app)

	authApi.AddAuthToApp(app)

	ui := app.Group("")
	ui.Use(authApi.GetAuthRedirectMiddleware())
	api.AddUIRoutes(ui, webhookStore)

	webhook := app.Group("webhook")
	api.AddWebhookRoutes(webhook, webhookStore, *publisher)

	configRest := app.Group("/api/v1/config")
	configRest.Use(authApi.GetAuthMiddleware())
	api.AddConfigRestRoutes(configRest, webhookStore)

	app.Listen(":3000")
}

func loadEnvironment(jiraUrl *string, jiraApiToken *string, tokenSignature *string) {
	godotenv.Load(".env")

	*jiraUrl = os.Getenv("JIRA_URL")
	*jiraApiToken = os.Getenv("JIRA_API_TOKEN")
	*tokenSignature = os.Getenv("TOKEN_SIGNATURE")

	log.Debugf("Jira URL: %s, Jira Token: %s", *jiraUrl, *jiraApiToken)
	log.Debugf("Token Signature: %s", *jiraUrl, *jiraApiToken)
}
