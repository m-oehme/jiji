package jira

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

func TestMapIssue_BasicFields(t *testing.T) {
	created := models.DateTimeScheme(time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC))
	updated := models.DateTimeScheme(time.Date(2026, 3, 18, 14, 30, 0, 0, time.UTC))

	iss := &models.IssueScheme{
		Key: "TEST-42",
		Fields: &models.IssueFieldsScheme{
			Summary: "Fix the login bug",
			Status:  &models.StatusScheme{Name: "In Progress", ID: "3"},
			Priority: &models.PriorityScheme{Name: "High"},
			IssueType: &models.IssueTypeScheme{Name: "Bug"},
			Assignee: &models.UserScheme{DisplayName: "Alice"},
			Reporter: &models.UserScheme{DisplayName: "Bob"},
			Labels:   []string{"backend", "auth"},
			Created:  &created,
			Updated:  &updated,
		},
	}

	issue := mapIssue(iss)

	if issue.Key != "TEST-42" {
		t.Errorf("Key: got %q, want TEST-42", issue.Key)
	}
	if issue.Summary != "Fix the login bug" {
		t.Errorf("Summary: got %q", issue.Summary)
	}
	if issue.Status != "In Progress" {
		t.Errorf("Status: got %q", issue.Status)
	}
	if issue.StatusID != "3" {
		t.Errorf("StatusID: got %q", issue.StatusID)
	}
	if issue.Priority != "High" {
		t.Errorf("Priority: got %q", issue.Priority)
	}
	if issue.Type != "Bug" {
		t.Errorf("Type: got %q", issue.Type)
	}
	if issue.Assignee != "Alice" {
		t.Errorf("Assignee: got %q", issue.Assignee)
	}
	if issue.Reporter != "Bob" {
		t.Errorf("Reporter: got %q", issue.Reporter)
	}
	if len(issue.Labels) != 2 || issue.Labels[0] != "backend" {
		t.Errorf("Labels: got %v", issue.Labels)
	}
	if issue.Created.Year() != 2026 {
		t.Errorf("Created: got %v", issue.Created)
	}
	if issue.Updated.Month() != time.March {
		t.Errorf("Updated: got %v", issue.Updated)
	}
}

func TestMapIssue_ADFDescription(t *testing.T) {
	adfNode := &models.CommentNodeScheme{
		Version: 1,
		Type:    "doc",
		Content: []*models.CommentNodeScheme{
			{
				Type: "paragraph",
				Content: []*models.CommentNodeScheme{
					{Type: "text", Text: "Hello from Jira"},
				},
			},
		},
	}

	iss := &models.IssueScheme{
		Key: "ADF-1",
		Fields: &models.IssueFieldsScheme{
			Summary:     "ADF test",
			Description: adfNode,
		},
	}

	issue := mapIssue(iss)

	if issue.Description == nil {
		t.Fatal("expected Description to be non-nil json.RawMessage")
	}

	var doc map[string]any
	if err := json.Unmarshal(issue.Description, &doc); err != nil {
		t.Fatalf("failed to unmarshal ADF: %v", err)
	}
	if doc["type"] != "doc" {
		t.Errorf("expected ADF type 'doc', got %v", doc["type"])
	}
}

func TestMapIssue_NilFields(t *testing.T) {
	iss := &models.IssueScheme{Key: "EMPTY-1"}
	issue := mapIssue(iss)

	if issue.Key != "EMPTY-1" {
		t.Errorf("Key: got %q", issue.Key)
	}
	if issue.Summary != "" {
		t.Errorf("Summary should be empty, got %q", issue.Summary)
	}
}

func TestMapIssue_Links(t *testing.T) {
	iss := &models.IssueScheme{
		Key: "LINK-1",
		Fields: &models.IssueFieldsScheme{
			Summary: "Issue with links",
			IssueLinks: []*models.IssueLinkScheme{
				{
					Type: &models.LinkTypeScheme{Outward: "blocks"},
					OutwardIssue: &models.LinkedIssueScheme{
						Key: "LINK-2",
						Fields: &models.IssueLinkFieldsScheme{
							Summary: "Blocked issue",
							Status:  &models.StatusScheme{Name: "Open"},
						},
					},
				},
			},
			Subtasks: []*models.IssueScheme{
				{
					Key: "LINK-3",
					Fields: &models.IssueFieldsScheme{
						Summary: "Subtask",
						Status:  &models.StatusScheme{Name: "Done"},
					},
				},
			},
		},
	}

	issue := mapIssue(iss)

	if len(issue.Links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(issue.Links))
	}
	if issue.Links[0].Key != "LINK-2" {
		t.Errorf("link key: got %q", issue.Links[0].Key)
	}
	if issue.Links[0].LinkType != "blocks" {
		t.Errorf("link type: got %q", issue.Links[0].LinkType)
	}

	if len(issue.Subtasks) != 1 {
		t.Fatalf("expected 1 subtask, got %d", len(issue.Subtasks))
	}
	if issue.Subtasks[0].Key != "LINK-3" {
		t.Errorf("subtask key: got %q", issue.Subtasks[0].Key)
	}
}

func TestMapComment(t *testing.T) {
	c := &models.IssueCommentScheme{
		ID:      "10001",
		Author:  &models.UserScheme{DisplayName: "Charlie"},
		Body: &models.CommentNodeScheme{
			Version: 1,
			Type:    "doc",
			Content: []*models.CommentNodeScheme{
				{Type: "paragraph", Content: []*models.CommentNodeScheme{
					{Type: "text", Text: "Nice work!"},
				}},
			},
		},
		Created: "2026-03-15T10:00:00.000+0000",
		Updated: "2026-03-15T10:05:00.000+0000",
	}

	comment := mapComment(c)

	if comment.ID != "10001" {
		t.Errorf("ID: got %q", comment.ID)
	}
	if comment.Author != "Charlie" {
		t.Errorf("Author: got %q", comment.Author)
	}
	if comment.Body == nil {
		t.Fatal("expected Body to be non-nil")
	}
	if comment.Created.IsZero() {
		t.Error("Created should not be zero")
	}
}

func TestParseJiraTime(t *testing.T) {
	tests := []struct {
		input string
		year  int
	}{
		{"2026-03-15T10:00:00.000+0000", 2026},
		{"2026-03-15T10:00:00.000+0200", 2026},
		{"2026-03-15T10:00:00Z", 2026},
		{"", 1},   // zero time
		{"invalid", 1},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJiraTime(tt.input)
			if result.Year() != tt.year {
				t.Errorf("parseJiraTime(%q) year = %d, want %d", tt.input, result.Year(), tt.year)
			}
		})
	}
}

func TestWrapError_NilResponse(t *testing.T) {
	err := wrapError(nil, fmt.Errorf("something failed"))
	if err == nil {
		t.Error("expected error")
	}
}
