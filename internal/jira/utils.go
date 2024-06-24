package jira

import "encoding/base64"

// According to the JIRA v2 API documentation, the JIRA API token should be converted to base64.
func ConvertJiraToken(jiraToken string) string {
	return base64.StdEncoding.EncodeToString([]byte(jiraToken))
}
