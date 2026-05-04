package jira

import (
	"testing"

	"github.com/m-oehme/jiji/internal/config"
)

func TestJira_IssueURL(t *testing.T) {
	tests := []struct {
		name, host, key, want string
	}{
		{"plain host", "https://co.atlassian.net", "PTECH-1", "https://co.atlassian.net/browse/PTECH-1"},
		{"trailing slash", "https://co.atlassian.net/", "PTECH-1", "https://co.atlassian.net/browse/PTECH-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Issue{
				Key: tt.key,
			}
			config := config.Config{
				Jira: config.JiraConnection{
					Host: tt.host,
				},
			}

			got := m.IssueURL(config)
			if got != tt.want {
				t.Fatalf("For Test '%s', expected: %s, got: %s", tt.name, tt.want, got)
			}
		})
	}
}
