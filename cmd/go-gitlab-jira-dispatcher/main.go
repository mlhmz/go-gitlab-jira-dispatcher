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

	loadEnvironment(&jiraUrl, &jiraApiToken)

	db := store.NewDatabase()
	var webhookStore store.WebhookConfigStore = sqlite.NewSqliteWebhookStore(db)

	publisher := gitlab.NewPublisher()
	publisher.Register(jira.NewJiraListener(jirav2.NewRestClient(jiraUrl, jiraApiToken)))

	signature := []byte("test-signature")
	var token auth.Token = auth.NewJwtToken(&signature, time.Hour*1)

	authMiddleware := api.AuthMiddleware(token)

	api.AddAuthToApp(app, db, token)

	ui := app.Group("")
	ui.Use(authMiddleware)
	api.AddUIRoutes(ui, webhookStore)

	webhook := app.Group("webhook")
	api.AddWebhookRoutes(webhook, webhookStore, *publisher)

	configRest := app.Group("/api/v1/config")
	configRest.Use(authMiddleware)
	api.AddConfigRestRoutes(configRest, webhookStore)

	app.Listen(":3000")
}

func loadEnvironment(jiraUrl *string, jiraApiToken *string) {
	godotenv.Load(".env")

	*jiraUrl = os.Getenv("JIRA_URL")
	*jiraApiToken = os.Getenv("JIRA_API_TOKEN")

	log.Debugf("Jira URL: %s, Jira Token: %s", *jiraUrl, *jiraApiToken)
}
