package jira

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
)

type JiraListener struct {
	client Client
}

func NewJiraListener(client Client) *JiraListener {
	return &JiraListener{
		client: client,
	}
}

func (listener *JiraListener) Accept(event *dispatcher.Event) {
	log.Infof("Jira listener accepted the event for the ticket '%s' with the status '%d' and the reviewer email '%s'",
		event.TicketNumber, event.StatusID, event.ReviewerEmail)
	listener.client.TransitionIssue(event.TicketNumber, event.StatusID, event.ReviewerEmail)
}
