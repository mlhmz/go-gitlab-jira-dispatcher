package gitlab_test

import (
	"testing"

	"github.com/mlhmz/go-gitlab-jira-dispatcher/internal/gitlab"
)

var goodTitles = map[string]string{
	"PROJ-123":   "PROJ-123: Fixing bugs",
	"PROJE-1234": "PROJE-1234: Fixing bugs",
	"PROJ-1234":  "Merge reqeust for PROJ-1234",
	"NICE-42":    "NICE-42: Fix connection timeout",
	"NICE-425":   "NICE-425: Fix connection timeout with a really long title TEST-123",
}

func TestResolveJiraTicketFromTitle_FillsTicketNumber_OnTitleContainingTicketNumber(t *testing.T) {
	for expected, title := range goodTitles {
		ticketNumber := ""
		if err := gitlab.ResolveJiraTicketFromTitle(title, &ticketNumber); err != nil {
			t.Errorf("ResolveJiraTicketFromTitle failed: %v", err)
		}
		if ticketNumber != expected {
			t.Errorf("ResolveJiraTicketFromTitle failed: expected %s, got %s", expected, ticketNumber)
		}
	}
}

func TestResolveJiraTicketFromTitle_ReturnsError_OnTitleNotContainingTicketNumber(t *testing.T) {
	ticketNumber := ""
	if err := gitlab.ResolveJiraTicketFromTitle("Update README.md", &ticketNumber); err == nil {
		t.Error("ResolveJiraTicketFromTitle failed: expected an error, got nil")
	}
}
