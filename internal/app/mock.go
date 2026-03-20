package app

import (
	"encoding/json"
	"time"

	"github.com/m-oehme/jiji/internal/jira"
)

// mockADF creates a simple ADF document with the given text.
func mockADF(text string) json.RawMessage {
	doc := map[string]any{
		"version": 1,
		"type":    "doc",
		"content": []map[string]any{
			{
				"type": "paragraph",
				"content": []map[string]any{
					{"type": "text", "text": text},
				},
			},
		},
	}
	b, _ := json.Marshal(doc)
	return b
}

// mockIssues returns test data for Phase 2 UI development.
func mockIssues() []jira.Issue {
	now := time.Now()
	return []jira.Issue{
		{
			Key: "PROJ-101", Summary: "Fix login timeout on mobile", Status: "In Progress", Priority: "High",
			Type: "Bug", Assignee: "Max", Reporter: "Sarah", Sprint: "Sprint 42",
			Labels: []string{"mobile", "auth"}, Created: now.Add(-72 * time.Hour), Updated: now.Add(-1 * time.Hour),
			Description: mockADF("Users experience a timeout when attempting to log in on mobile devices. The session token expires before the OAuth flow completes."),
			Subtasks: []jira.IssueLink{
				{Key: "PROJ-102", Summary: "Increase token TTL", Status: "Done", LinkType: "subtask"},
				{Key: "PROJ-103", Summary: "Add retry logic", Status: "To Do", LinkType: "subtask"},
			},
		},
		{
			Key: "PROJ-104", Summary: "Add caching layer for dashboard", Status: "To Do", Priority: "Medium",
			Type: "Story", Assignee: "Max", Reporter: "Alex", Sprint: "Sprint 42",
			Created: now.Add(-48 * time.Hour), Updated: now.Add(-6 * time.Hour),
			Description: mockADF("Dashboard queries are slow. Add Redis caching for the most frequent queries."),
		},
		{
			Key: "PROJ-105", Summary: "Update onboarding flow", Status: "In Review", Priority: "Medium",
			Type: "Story", Assignee: "Lisa", Reporter: "Max", Sprint: "Sprint 42",
			Labels: []string{"ux"}, Created: now.Add(-96 * time.Hour), Updated: now.Add(-2 * time.Hour),
			Description: mockADF("Redesign the onboarding wizard to reduce drop-off. New designs in Figma."),
		},
		{
			Key: "PROJ-106", Summary: "Migrate to PostgreSQL 16", Status: "To Do", Priority: "Low",
			Type: "Task", Assignee: "Chen", Reporter: "Sarah",
			Created: now.Add(-120 * time.Hour), Updated: now.Add(-24 * time.Hour),
			Description: mockADF("Plan and execute migration from PostgreSQL 14 to 16. Coordinate with DBA team."),
		},
		{
			Key: "PROJ-107", Summary: "API rate limiting", Status: "In Progress", Priority: "High",
			Type: "Story", Assignee: "Max", Reporter: "Alex", Sprint: "Sprint 42",
			Labels: []string{"api", "security"}, Created: now.Add(-36 * time.Hour), Updated: now.Add(-3 * time.Hour),
			Description: mockADF("Implement rate limiting on public API endpoints. Use token bucket algorithm."),
			Links: []jira.IssueLink{
				{Key: "PROJ-108", Summary: "Load test rate limiter", Status: "To Do", LinkType: "is blocked by"},
			},
		},
		{
			Key: "PROJ-108", Summary: "Load test rate limiter", Status: "To Do", Priority: "Medium",
			Type: "Task", Assignee: "Lisa", Reporter: "Max",
			Created: now.Add(-36 * time.Hour), Updated: now.Add(-36 * time.Hour),
			Description: mockADF("Set up k6 load tests for the new rate limiting endpoints."),
		},
		{
			Key: "PROJ-109", Summary: "Dark mode support", Status: "Done", Priority: "Low",
			Type: "Story", Assignee: "Chen", Reporter: "Lisa", Sprint: "Sprint 41",
			Labels: []string{"ux", "theme"}, Created: now.Add(-168 * time.Hour), Updated: now.Add(-48 * time.Hour),
			Description: mockADF("Add dark mode theme toggle. CSS variables already support it, just wire up the toggle."),
		},
		{
			Key: "PROJ-110", Summary: "Fix broken CSV export", Status: "In Progress", Priority: "Critical",
			Type: "Bug", Assignee: "Sarah", Reporter: "Chen",
			Created: now.Add(-12 * time.Hour), Updated: now.Add(-1 * time.Hour),
			Description: mockADF("CSV export truncates fields containing commas. Need to properly escape quoted fields."),
		},
		{
			Key: "PROJ-111", Summary: "Add Prometheus metrics", Status: "To Do", Priority: "Medium",
			Type: "Task", Assignee: "Alex", Reporter: "Max", Sprint: "Sprint 43",
			Created: now.Add(-24 * time.Hour), Updated: now.Add(-24 * time.Hour),
			Description: mockADF("Instrument key service methods with Prometheus counters and histograms."),
		},
		{
			Key: "PROJ-112", Summary: "Refactor user service", Status: "In Review", Priority: "Low",
			Type: "Task", Assignee: "Max", Reporter: "Sarah", Sprint: "Sprint 42",
			Created: now.Add(-144 * time.Hour), Updated: now.Add(-12 * time.Hour),
			Description: mockADF("Extract permission checks into middleware. The user service has grown too large."),
		},
		{
			Key: "PROJ-113", Summary: "Upgrade Go to 1.23", Status: "Done", Priority: "Low",
			Type: "Task", Assignee: "Chen", Reporter: "Alex",
			Created: now.Add(-200 * time.Hour), Updated: now.Add(-72 * time.Hour),
			Description: mockADF("Upgrade from Go 1.22 to 1.23. Check for breaking changes in dependencies."),
		},
		{
			Key: "PROJ-114", Summary: "Email notification throttling", Status: "To Do", Priority: "High",
			Type: "Story", Assignee: "Lisa", Reporter: "Sarah", Sprint: "Sprint 43",
			Created: now.Add(-6 * time.Hour), Updated: now.Add(-6 * time.Hour),
			Description: mockADF("Users complain about too many email notifications. Implement digest mode."),
		},
		{
			Key: "PROJ-115", Summary: "Accessibility audit fixes", Status: "To Do", Priority: "Medium",
			Type: "Story", Assignee: "Sarah", Reporter: "Lisa", Sprint: "Sprint 43",
			Labels: []string{"a11y", "ux"}, Created: now.Add(-48 * time.Hour), Updated: now.Add(-48 * time.Hour),
			Description: mockADF("Fix WCAG 2.1 AA violations found in the accessibility audit."),
		},
		{
			Key: "PROJ-116", Summary: "S3 bucket lifecycle policies", Status: "In Progress", Priority: "Low",
			Type: "Task", Assignee: "Alex", Reporter: "Chen",
			Created: now.Add(-72 * time.Hour), Updated: now.Add(-8 * time.Hour),
			Description: mockADF("Configure lifecycle policies to move old objects to Glacier after 90 days."),
		},
		{
			Key: "PROJ-117", Summary: "WebSocket reconnection logic", Status: "To Do", Priority: "High",
			Type: "Bug", Assignee: "Max", Reporter: "Lisa", Sprint: "Sprint 42",
			Labels: []string{"realtime"}, Created: now.Add(-24 * time.Hour), Updated: now.Add(-2 * time.Hour),
			Description: mockADF("WebSocket connections drop silently on network change. Need exponential backoff reconnection."),
		},
	}
}

// mockComments returns test comments for Phase 2 UI development.
func mockComments() []jira.Comment {
	now := time.Now()
	return []jira.Comment{
		{
			ID: "10001", Author: "Sarah",
			Body:    mockADF("I can reproduce this on iOS 17. The timeout happens after about 30 seconds."),
			Created: now.Add(-48 * time.Hour), Updated: now.Add(-48 * time.Hour),
		},
		{
			ID: "10002", Author: "Max",
			Body:    mockADF("Root cause found: the OAuth redirect takes too long on cellular connections. Increasing the token TTL to 5 minutes should fix it."),
			Created: now.Add(-24 * time.Hour), Updated: now.Add(-24 * time.Hour),
		},
		{
			ID: "10003", Author: "Alex",
			Body:    mockADF("Should we also add a retry mechanism? In case the first attempt fails."),
			Created: now.Add(-12 * time.Hour), Updated: now.Add(-12 * time.Hour),
		},
	}
}
