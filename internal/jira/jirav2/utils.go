package jirav2

import "encoding/base64"

// According to the JIRA v2 API documentation, the JIRA API token should be converted to base64.
func convertJiraToken(jiraToken string) string {
	return base64.StdEncoding.EncodeToString([]byte(jiraToken))
}
