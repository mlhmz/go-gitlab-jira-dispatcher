package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"mlhmz.xyz/go-gitlab-jira-dispatcher/internal/gitlab"
)

func main() {
	app := fiber.New()

	app.Post("/webhook", func(c *fiber.Ctx) error {
		var event *gitlab.MergeRequestEvent

		log.Info(string(c.Body()[:]))

		if err := c.BodyParser(&event); err != nil {
			log.Error("Failed to parse webhook event", "error", err)
			return c.SendStatus(400)
		}

		log.Info("Received webhook event with the action ", event.ObjectAttributes.Action, " and the name '", event.ObjectAttributes.Title, "'")

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
