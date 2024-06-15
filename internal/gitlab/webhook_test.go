package gitlab

import (
	"testing"

	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/dispatcher"
	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/store"
)

var transitions = store.Transitions{
	ReadyForReview:  1,
	InReview:        2,
	ReviewOK:        3,
	ReviewNotOK:     4,
	DevelopmentDone: 5,
}

func TestWebhookPublisher_Register(t *testing.T) {
	publisher := NewPublisher(&transitions)
	listener := &mockListener{}

	publisher.Register(listener)

	if len(publisher.listeners) != 1 {
		t.Errorf("Expected 1 listener, got %d", len(publisher.listeners))
	}
}

type SuccessTestCase struct {
	input    MergeRequestEvent
	expected dispatcher.Event
}

var successTestCases = []SuccessTestCase{
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "open",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.ReadyForReview,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "reopen",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.ReadyForReview,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "merge",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.DevelopmentDone,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "approved",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.ReviewOK,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "unapproved",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.InReview,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "close",
			},
		},
		expected: dispatcher.Event{
			TicketNumber: "TEST-1000",
			StatusID:     transitions.ReviewNotOK,
		},
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "update",
			},
			Changes: Changes{
				Reviewers: UserChange{
					Previous: []User{},
					Current: []User{
						{
							Email: "test@example.org",
						},
					},
				},
			},
		},
		expected: dispatcher.Event{
			TicketNumber:  "TEST-1000",
			StatusID:      transitions.InReview,
			ReviewerEmail: "test@example.org",
		},
	},
}

func TestWebhookPublisher_ProcessWebhook(t *testing.T) {
	for _, testCase := range successTestCases {
		publisher := NewPublisher(&transitions)
		listener := &mockListener{}
		dispatcherEvent := &dispatcher.Event{}

		publisher.Register(listener)

		err := publisher.ProcessWebhook(&testCase.input, dispatcherEvent)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if listener.event == nil {
			t.Errorf("Expected listener event to be not nil")
		}

		if listener.event.TicketNumber != testCase.expected.TicketNumber {
			t.Errorf("Expected ticket number to be '%s', got '%s'", testCase.expected.TicketNumber, listener.event.TicketNumber)
		}

		if listener.event.StatusID != testCase.expected.StatusID {
			t.Errorf("Expected status ID to be '%d', got '%d'", testCase.expected.StatusID, listener.event.StatusID)
		}

		if listener.event.ReviewerEmail != testCase.expected.ReviewerEmail {
			t.Errorf("Expected reviewer email to be '%s', got '%s'", testCase.expected.ReviewerEmail, listener.event.ReviewerEmail)
		}
	}
}

type ErrorTestCase struct {
	input         MergeRequestEvent
	expectedError string
}

var errorTestCases = []ErrorTestCase{
	{
		input: MergeRequestEvent{
			ObjectKind: "wrong_kind",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "open",
			},
		},
		expectedError: "'wrong_kind' is a invalid event type. Only 'merge_request' events are supported",
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "TEST-1000: Merge Request Title",
				Action: "wrongaction",
			},
		},
		expectedError: "no action found for the event 'wrongaction'",
	},
	{
		input: MergeRequestEvent{
			ObjectKind: "merge_request",
			ObjectAttributes: ObjectAttributes{
				Title:  "Merge Request Title",
				Action: "open",
			},
		},
		expectedError: "no Jira ticket found in the title 'Merge Request Title'",
	},
}

func TestWebhookPublisher_ProcessWebhook_ReturnsError_OnErrorCases(t *testing.T) {
	for _, testCase := range errorTestCases {
		publisher := NewPublisher(&transitions)
		listener := &mockListener{}
		dispatcherEvent := &dispatcher.Event{}

		publisher.Register(listener)

		err := publisher.ProcessWebhook(&testCase.input, dispatcherEvent)

		if err == nil {
			t.Errorf("Expected an error, got %v", err)
		}

		if err.Error() != testCase.expectedError {
			t.Errorf("Expected error message to be '%s', got '%s'", testCase.expectedError, err.Error())
		}

		if listener.event != nil {
			t.Errorf("Expected listener event to be nil, got status %v", listener.event.StatusID)
		}
	}
}

type mockListener struct {
	event *dispatcher.Event
}

func (m *mockListener) Accept(event *dispatcher.Event) {
	m.event = event
}
