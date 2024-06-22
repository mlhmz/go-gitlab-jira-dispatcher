package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
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

	app.Get("", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.Render("index", fiber.Map{
			"Configs": configs,
		}, "layouts/main")
	})

	app.Get("/config-list", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.Render("config-list", fiber.Map{
			"Configs": configs,
		})
	})

	app.Post("/config", func(c *fiber.Ctx) error {
		var submission store.WebhookConfigSubmission

		if err := c.BodyParser(&submission); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		config := store.MapSubmission(&submission)
		if err := webhookStore.CreateWebhookConfig(config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		c.Response().Header.Add("HX-Trigger", "reloadConfig")
		return c.Render("config", fiber.Map{
			"Config": config,
		})
	})

	app.Get("/config/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)

		if err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		var config store.WebhookConfig
		if err := webhookStore.GetWebhookConfig(uuid, &config); err != nil {
			log.Warn(err)
			return c.Status(404).SendString(err.Error())
		}

		return c.Render("config", fiber.Map{
			"Config": config,
		})
	})

	app.Get("/create", func(c *fiber.Ctx) error {
		return c.Render("create", fiber.Map{})
	})

	app.Post("/webhook/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)

		if err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		var config store.WebhookConfig
		if err := webhookStore.GetWebhookConfig(uuid, &config); err != nil {
			log.Warn(err)
			return c.Status(404).SendString(err.Error())
		}

		var event gitlab.MergeRequestEvent

		if err := c.BodyParser(&event); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook event error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		result := dispatcher.Event{}
		if err := publisher.ProcessWebhook(&event, &result, &config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(result)
	})

	configApi := app.Group("/api/v1/config")

	configApi.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)

		if err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		var config store.WebhookConfig
		if err := webhookStore.GetWebhookConfig(uuid, &config); err != nil {
			log.Warn(err)
			return c.Status(404).SendString(err.Error())
		}

		return c.JSON(config)
	})

	configApi.Get("", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(configs)
	})

	configApi.Post("", func(c *fiber.Ctx) error {
		var config store.WebhookConfig

		if err := c.BodyParser(&config); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		if err := webhookStore.CreateWebhookConfig(&config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(config)
	})

	configApi.Put("", func(c *fiber.Ctx) error {
		var config store.WebhookConfig
		if err := c.BodyParser(&config); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		if err := webhookStore.UpdateWebhookConfig(&config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.JSON(config)
	})

	configApi.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		uuid, err := uuid.Parse(id)

		if err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		if err := webhookStore.DeleteWebhookConfig(uuid); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		c.Response().Header.Add("HX-Trigger", "reloadConfig,deleteConfig")
		return c.SendStatus(204)
	})

	app.Listen(":3000")
}

func loadEnvironment(jiraUrl *string, jiraApiToken *string) {
	godotenv.Load(".env")

	*jiraUrl = os.Getenv("JIRA_URL")
	*jiraApiToken = os.Getenv("JIRA_API_TOKEN")

	log.Debugf("Jira URL: %s, Jira Token: %s", *jiraUrl, *jiraApiToken)
}
