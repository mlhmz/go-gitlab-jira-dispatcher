package jira

type Client interface {
	TransitionIssue(ticketNumber string, statusID int, reviewerEmail string) error
}
