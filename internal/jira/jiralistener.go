package jira

import (
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
	listener.client.TransitionIssue(event.TicketNumber, event.StatusID, event.ReviewerEmail)
}
