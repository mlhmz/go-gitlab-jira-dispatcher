package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/jira"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/jira/jirav2"
)

func main() {
	app := fiber.New()

	var jiraUrl string
	var jiraApiToken string

	loadEnvironment(&jiraUrl, &jiraApiToken)

	publisher := gitlab.NewPublisher()
	publisher.Register(jira.NewJiraListener(jirav2.NewRestClient(jiraUrl, jiraApiToken)))

	app.Post("/webhook", func(c *fiber.Ctx) error {
		var event gitlab.MergeRequestEvent

		if err := c.BodyParser(&event); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook event error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		result := dispatcher.Event{}
		if err := publisher.ProcessWebhook(&event, &result); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(result)
	})

	app.Listen(":3000")
}

func loadEnvironment(jiraUrl *string, jiraApiToken *string) {
	godotenv.Load(".env")

	*jiraUrl = os.Getenv("JIRA_URL")
	*jiraApiToken = os.Getenv("JIRA_API_TOKEN")

	log.Debugf("Jira URL: %s, Jira Token: %s", *jiraUrl, *jiraApiToken)
}
