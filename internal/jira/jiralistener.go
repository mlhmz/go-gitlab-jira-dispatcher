package jira

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
)

type JiraListener struct {
}

func NewJiraListener() *JiraListener {
	return &JiraListener{}
}

func (j *JiraListener) Accept(event *dispatcher.Event) {
	log.Infof("Jira listener accepted the event for the ticket '%s' with the status '%d' and the reviewer email '%s'",
		event.TicketNumber, event.StatusID, event.ReviewerEmail)
}
