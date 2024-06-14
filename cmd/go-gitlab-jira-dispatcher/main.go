package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
)

func main() {
	app := fiber.New()

	app.Post("/webhook", func(c *fiber.Ctx) error {
		var event *gitlab.MergeRequestEvent

		if err := c.BodyParser(&event); err != nil {
			log.Error("Failed to parse webhook event", "error", err)
			return c.SendStatus(400)
		}

		var ticketNumber string
		if err := gitlab.ResolveJiraTicketFromTitle(event.ObjectAttributes.Title, &ticketNumber); err != nil {
			log.Warn(err)
			return c.SendStatus(400)
		}

		action := gitlab.NewAction(event.ObjectAttributes.Action)

		result := action.Execute(ticketNumber, event)

		if result != nil {
			log.Infof("Dispatched event for the ticket '%s' with the status '%s' and the reviewer email '%s'",
				result.TicketNumber, result.Status, result.ReviewerEmail)
		}

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
