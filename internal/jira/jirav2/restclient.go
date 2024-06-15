package jirav2

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var Prefix = "/rest/api/2"

type RestClient struct {
	url   string
	token string
}

func NewRestClient(url string, token string) *RestClient {
	return &RestClient{
		url:   url,
		token: token,
	}
}

func (r *RestClient) TransitionIssue(ticketNumber string, statusID int, reviewerEmail string) error {
	url := r.url + Prefix + fmt.Sprintf("/issue/%s/transitions", ticketNumber)

	log.Debugf("Requesting to jira with the following url: %s", url)

	body := &TransitionPayload{
		Transition: Transition{
			ID: fmt.Sprint(statusID),
		},
	}
	agent := fiber.Post(url)

	agent.JSON(body)
	agent.Set("Content-Type", "application/json")
	agent.Set("Authorization", "Basic "+convertJiraToken(r.token))

	status, res, err := agent.Bytes()

	if len(err) > 0 {
		return err[0]
	}

	log.Infof("Jira v2 Restclient: Transitioned issue '%s' to status '%d' with status code '%d' and the message '%s'",
		ticketNumber, statusID, status, string(res[:]))

	return nil
}
