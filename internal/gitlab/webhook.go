package gitlab

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
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

func (publisher *WebhookPublisher) Notify(event *dispatcher.Event) {
	for _, listener := range publisher.listeners {
		listener.Accept(event)
	}
}

func (publisher *WebhookPublisher) ProcessWebhook(mrEvent *MergeRequestEvent, dispatcherEvent *dispatcher.Event) error {
	if mrEvent.ObjectKind != "merge_request" {
		return fmt.Errorf("'%s' is a invalid event type. Only 'merge_request' events are supported",
			mrEvent.ObjectKind)
	}

	var ticketNumber string
	if err := ResolveJiraTicketFromTitle(mrEvent.ObjectAttributes.Title, &ticketNumber); err != nil {
		return err
	}

	action := NewAction(mrEvent.ObjectAttributes.Action)

	if action == nil {
		return fmt.Errorf("no action found for the event '%s'", mrEvent.ObjectAttributes.Action)
	}

	if actionResult := action.Execute(ticketNumber, mrEvent); actionResult != nil {
		*dispatcherEvent = *actionResult
		log.Infof("Dispatched event for the ticket '%s' with the status '%s' and the reviewer email '%s'",
			dispatcherEvent.TicketNumber, dispatcherEvent.StatusID, dispatcherEvent.ReviewerEmail)
		publisher.Notify(dispatcherEvent)
	}
	return nil
}
