package gitlab

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

type WebhookPublisher struct {
	listeners []dispatcher.Listener
}

func NewPublisher() *WebhookPublisher {
	return &WebhookPublisher{}
}

func (publisher *WebhookPublisher) Register(listener dispatcher.Listener) {
	publisher.listeners = append(publisher.listeners, listener)
}

func (publisher *WebhookPublisher) ProcessWebhook(mrEvent *MergeRequestEvent, dispatcherEvent *dispatcher.Event, config *store.WebhookConfig) error {
	if mrEvent.ObjectKind != "merge_request" {
		return fmt.Errorf("'%s' is a invalid event type. Only 'merge_request' events are supported",
			mrEvent.ObjectKind)
	}

	var ticketNumber string
	if err := ResolveJiraTicketFromTitle(mrEvent.ObjectAttributes.Title, &ticketNumber); err != nil {
		return err
	}

	projects := strings.Split(config.Projects, ",")
	if !IsTicketNumberWhitelisted(&projects, &ticketNumber) {
		return fmt.Errorf("'%s' is not allowed for this config, only [%s] are allowed", ticketNumber, projects)
	}

	action := NewAction(mrEvent.ObjectAttributes.Action)

	if action == nil {
		return fmt.Errorf("no action found for the event '%s'", mrEvent.ObjectAttributes.Action)
	}

	if actionResult := action.Execute(&ticketNumber, mrEvent, config); actionResult != nil {
		*dispatcherEvent = *actionResult
		log.Infof("Dispatched event for the ticket '%s' with the status '%d' and the reviewer email '%s'",
			dispatcherEvent.TicketNumber, dispatcherEvent.StatusID, dispatcherEvent.ReviewerEmail)
		publisher.notify(dispatcherEvent)
	}
	return nil
}

func (publisher *WebhookPublisher) notify(event *dispatcher.Event) {
	for _, listener := range publisher.listeners {
		listener.Accept(event)
	}
}
