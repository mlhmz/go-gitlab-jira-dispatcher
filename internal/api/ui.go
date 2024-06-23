package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

func AddUIRoutes(router fiber.Router, webhookStore store.WebhookConfigStore) {
	router.Get("", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.Render("index", fiber.Map{
			"Configs": configs,
		}, "layouts/main")
	})

	router.Get("/config-list", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := webhookStore.GetAllWebhookConfigs(&configs); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		return c.Render("config-list", fiber.Map{
			"Configs": configs,
		})
	})

	router.Post("/config", func(c *fiber.Ctx) error {
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

	router.Get("/config/:id", func(c *fiber.Ctx) error {
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

	router.Get("/create", func(c *fiber.Ctx) error {
		return c.Render("create", fiber.Map{})
	})

	// In order to remove an element with htmx on delete return 200
	// with an empty body
	router.Delete("/config/:id", func(c *fiber.Ctx) error {
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
		return c.SendString("")
	})
}
