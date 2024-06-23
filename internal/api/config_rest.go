package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

func AddConfigRestRoutes(router fiber.Router, webhookStore store.WebhookConfigStore) {
	router.Get("/:id", func(c *fiber.Ctx) error {
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

	router.Get("", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
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

		c.Response().Header.Add("HX-Trigger", "reloadConfig")
		return c.SendStatus(204)
	})
}
