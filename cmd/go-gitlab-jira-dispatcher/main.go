package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/jira"
)

func main() {
	app := fiber.New(fiber.Config{})

	var jiraUrl string
	var jiraApiToken string

	loadEnvironment(&jiraUrl, &jiraApiToken)

	app.Post("/webhook/:readyForReview<int\\>/:inReview<int\\>/:reviewOK<int\\>/:reviewNotOK<int\\>/:developmentDone<int\\>", func(c *fiber.Ctx) error {
		config := map[string]int{
			"open":       parseInt(c.Params("readyForReview")),
			"reopen":     parseInt(c.Params("readyForReview")),
			"update":     parseInt(c.Params("inReview")),
			"unapproved": parseInt(c.Params("inReview")),
			"approved":   parseInt("reviewOK"),
			"close":      parseInt("reviewNotOK"),
			"merge":      parseInt("developmentDone"),
		}

		var event gitlab.MergeRequestEvent

		if err := c.BodyParser(&event); err != nil {
			formattedError := fmt.Errorf("failed to parse webhook event error: %s", err)
			log.Error(formattedError)
			return c.Status(500).SendString(formattedError.Error())
		}

		item, exists := config[event.ObjectAttributes.Action]

		if exists {
			if event.ObjectAttributes.Action == "update" && len(event.Changes.Reviewers.Current) < 1 {
				log.Info("the update event will not be forwarded to jira, because there was no reviewer change")
				return c.SendStatus(200)
			}

			var ticketNumber string
			url := jiraUrl + fmt.Sprintf("/rest/api/2/issue/%s/transitions", gitlab.ResolveJiraTicketFromTitle(event.ObjectAttributes.Title, &ticketNumber))

			log.Debugf("Requesting to jira with the following url: %s", url)

			body := &jira.TransitionPayload{
				Transition: jira.Transition{
					ID: fmt.Sprint(item),
				},
			}
			agent := fiber.Post(url)

			agent.JSON(body)
			agent.Set("Content-Type", "application/json")
			agent.Set("Authorization", "Basic "+jira.ConvertJiraToken(jiraApiToken))

			status, res, err := agent.Bytes()

			if len(err) > 0 {
				return c.Status(status).SendString(string(res[:]))
			}
			return c.JSON(string(res[:]))
		} else {
			eventError := fmt.Errorf("the event with the action '%s' was not found.", event.ObjectAttributes.Action)
			log.Error(eventError)
			return c.Status(404).SendString(eventError.Error())
		}
	})

	app.Listen(":3000")
}

func parseInt(value string) int {
	result, err := strconv.Atoi(value)

	if err == nil {
		return result
	} else {
		log.Errorf("there was an error while parsing the value %s - defaulting to 0", value)
		return 0
	}
}

func loadEnvironment(jiraUrl *string, jiraApiToken *string) {
	godotenv.Load(".env")

	*jiraUrl = os.Getenv("JIRA_URL")
	*jiraApiToken = os.Getenv("JIRA_API_TOKEN")

	log.Debugf("Jira URL: %s, Jira Token: %s", *jiraUrl, *jiraApiToken)
}
