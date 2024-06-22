package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

func AddWebhookRoutes(router fiber.Router, webhookStore store.WebhookConfigStore, publisher gitlab.WebhookPublisher) {
	router.Post("/webhook/:id", func(c *fiber.Ctx) error {
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
}
