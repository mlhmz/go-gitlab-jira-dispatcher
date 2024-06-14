package gitlab

import (
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
)

type Action interface {
	Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event
}

type OpenAction struct{}

func (a *OpenAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "Ready for Review",
	}
}

type ReopenAction struct{}

func (a *ReopenAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "Ready for Review",
	}
}

type UpdateAction struct{}

func (a *UpdateAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	if len(event.Changes.Reviewers.Previous) == 0 && len(event.Changes.Reviewers.Current) > 0 {
		return &dispatcher.Event{
			TicketNumber:  ticketNumber,
			Status:        "In Review",
			ReviewerEmail: event.Changes.Reviewers.Current[0].Email,
		}
	} else {
		return nil
	}
}

type MergeAction struct{}

func (a *MergeAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "Development done",
	}
}

type ApprovedAction struct{}

func (a *ApprovedAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "Review OK",
	}
}

type UnapprovedAction struct{}

func (a *UnapprovedAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "In Review",
	}
}

type CloseAction struct{}

func (a *CloseAction) Execute(ticketNumber string, event *MergeRequestEvent) *dispatcher.Event {
	return &dispatcher.Event{
		TicketNumber: ticketNumber,
		Status:       "Review not OK",
	}
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
