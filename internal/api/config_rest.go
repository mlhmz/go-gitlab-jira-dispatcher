package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

func AddConfigRestRoutes(router fiber.Router, webhookStore store.WebhookConfigStore) {
	commons := newCommons(&webhookStore)

	router.Get("/:id", func(c *fiber.Ctx) error {
		var config store.WebhookConfig
		if err := commons.getWebhookConfigByCtx(c, &config); err != nil {
			return err
		}
		return c.JSON(config)
	})

	router.Get("", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := commons.getAllWebhookConfigs(c, &configs); err != nil {
			return err
		}
		return c.JSON(configs)
	})

	router.Post("", func(c *fiber.Ctx) error {
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

	router.Put("", func(c *fiber.Ctx) error {
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

	router.Delete("/:id", func(c *fiber.Ctx) error {
		if err := commons.deleteWebhookConfigByCtx(c); err != nil {
			return err
		}
		return c.SendStatus(204)
	})
}
