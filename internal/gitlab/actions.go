package gitlab

import (
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

type Action interface {
	Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event
}

func NewAction(action string) Action {
	switch action {
	case "open":
		return &OpenAction{}
	case "reopen":
		return &ReopenAction{}
	case "update":
		return &UpdateAction{}
	case "merge":
		return &MergeAction{}
	case "approved":
		return &ApprovedAction{}
	case "unapproved":
		return &UnapprovedAction{}
	case "close":
		return &CloseAction{}
	default:
		return nil
	}
}

type OpenAction struct{}

func (a *OpenAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.ReadyForReview,
	}
}

type ReopenAction struct{}

func (a *ReopenAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.ReadyForReview,
	}
}

type UpdateAction struct{}

// GitLab will send a merge request event for every update that is triggered (e.g. changing the title)
// In order to detect if a reviewer is actually added to the merge request, we need to check first, if there
// was no reviewer before and if there is a reviewer now.
func (a *UpdateAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	if len(event.Changes.Reviewers.Previous) == 0 && len(event.Changes.Reviewers.Current) > 0 {
		return &dispatcher.Event{
			TicketNumber:  *ticketNumber,
			StatusID:      transitions.InReview,
			ReviewerEmail: event.Changes.Reviewers.Current[0].Email,
		}
	} else {
		return nil
	}
}

type MergeAction struct{}

func (a *MergeAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.DevelopmentDone,
	}
}

type ApprovedAction struct{}

func (a *ApprovedAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.ReviewOK,
	}
}

type UnapprovedAction struct{}

func (a *UnapprovedAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.InReview,
	}
}

type CloseAction struct{}

func (a *CloseAction) Execute(ticketNumber *string, event *MergeRequestEvent, transitions *store.WebhookConfig) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: *ticketNumber,
		StatusID:     transitions.ReviewNotOK,
	}
}
