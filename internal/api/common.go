package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

// as this application is supporting a hypermedia-driven api (ui with htmx) as well as
// an actual rest api, some fiber and crud connection code is duplicated,
// thats why a commons structure is created here.
type commons struct {
	webhookStore store.WebhookConfigStore
}

func newCommons(store *store.WebhookConfigStore) *commons {
	return &commons{
		webhookStore: *store,
	}
}

func (c *commons) getAllWebhookConfigs(ctx *fiber.Ctx, configs *[]store.WebhookConfig) error {
	if err := c.webhookStore.GetAllWebhookConfigs(configs); err != nil {
		log.Warn(err)
		return ctx.Status(400).SendString(err.Error())
	}
	return nil
}

func (c *commons) getWebhookConfigByCtx(ctx *fiber.Ctx, config *store.WebhookConfig) error {
	id := ctx.Params("id")
	uuid, err := uuid.Parse(id)

	if err != nil {
		log.Warn(err)
		return ctx.Status(400).SendString(err.Error())
	}

	if err := c.webhookStore.GetWebhookConfig(uuid, config); err != nil {
		log.Warn(err)
		return ctx.Status(404).SendString(err.Error())
	}
	return nil
}

func (c *commons) deleteWebhookConfigByCtx(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	uuid, err := uuid.Parse(id)

	if err != nil {
		log.Warn(err)
		return ctx.Status(400).SendString(err.Error())
	}

	if err := c.webhookStore.DeleteWebhookConfig(uuid); err != nil {
		log.Warn(err)
		return ctx.Status(400).SendString(err.Error())
	}
	return nil
}
