package gitlab

import (
	"fmt"
	"regexp"
	"strings"
)

func ResolveJiraTicketFromTitle(title string, ticketNumber *string) error {
	pattern := regexp.MustCompile(`([A-Z]+-\d+)`)
	matches := pattern.FindStringSubmatch(title)
	if len(matches) > 1 {
		*ticketNumber = matches[0]
		return nil
	} else {
		return fmt.Errorf("no Jira ticket found in the title '%s'", title)
	}
}

func IsTicketNumberWhitelisted(projects *[]string, ticketNumber *string) bool {
	project := strings.Split(*ticketNumber, "-")[0]
	for _, entry := range *projects {
		if entry == project {
			return true
		}
	}
	return false
}
