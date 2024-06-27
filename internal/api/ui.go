package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

// The UI of the application is hypermedia-driven with htmx.
// So instead of return json, actual, already filled, html is returned.
func AddUIRoutes(router fiber.Router, webhookStore store.WebhookConfigStore) {
	commons := newCommons(&webhookStore)

	router.Get("", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

	router.Get("/config-list", func(c *fiber.Ctx) error {
		var configs []store.WebhookConfig
		if err := commons.getAllWebhookConfigs(c, &configs); err != nil {
			return err
		}
		return c.Render("config-list", fiber.Map{
			"Configs": configs,
		}, "layouts/app")
	})

	// this post is a htmx specific post, as htmx is following the official html form conventions,
	// integers are submitted as strings and have to be parsed in the backend
	router.Post("/config", func(c *fiber.Ctx) error {
		var submission store.WebhookConfigSubmission

		if err := c.BodyParser(&submission); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		config := store.CreateFromSubmission(&submission)
		if err := webhookStore.CreateWebhookConfig(config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}

		c.Response().Header.Add("HX-Trigger", "reloadConfig")
		return c.Render("config", fiber.Map{
			"Config": config,
		})
	})

	router.Put("/:id", func(c *fiber.Ctx) error {
		var config store.WebhookConfig
		if err := commons.getWebhookConfigByCtx(c, &config); err != nil {
			return err
		}

		var submission store.WebhookConfigSubmission

		if err := c.BodyParser(&submission); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		config.UpdateFromSubmission(&submission)
		if err := webhookStore.UpdateWebhookConfig(&config); err != nil {
			log.Warn(err)
			return c.Status(400).SendString(err.Error())
		}
		return c.Render("config", fiber.Map{
			"Config":  config,
			"Message": fmt.Sprintf("The config with the id '%s' was successfully updated.", config.ID),
		})
	})

	router.Get("/config/:id", func(c *fiber.Ctx) error {
		var config store.WebhookConfig
		if err := commons.getWebhookConfigByCtx(c, &config); err != nil {
			return err
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
		if err := commons.deleteWebhookConfigByCtx(c); err != nil {
			return err
		}
		c.Response().Header.Add("HX-Trigger", "reloadConfig")
		return c.SendString("")
	})
}
