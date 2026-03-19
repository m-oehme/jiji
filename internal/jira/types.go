package jira

import (
	"encoding/json"
	"time"
)

// Issue represents a Jira issue with raw ADF for rich text fields.
type Issue struct {
	Key          string
	Summary      string
	Description  json.RawMessage // Raw ADF per ADR-011
	Status       string
	StatusID     string
	Priority     string
	Type         string
	Assignee     string
	Reporter     string
	Labels       []string
	Components   []string
	Sprint       string
	FixVersions  []string
	Created      time.Time
	Updated      time.Time
	Subtasks     []IssueLink
	Links        []IssueLink
	CustomFields map[string]any
}

// IssueLink represents a relationship to another issue.
type IssueLink struct {
	Key      string
	Summary  string
	Status   string
	LinkType string
}

// Comment represents a Jira comment with raw ADF body.
type Comment struct {
	ID      string
	Author  string
	Body    json.RawMessage // Raw ADF per ADR-011
	Created time.Time
	Updated time.Time
}

// Transition represents an available workflow transition.
type Transition struct {
	ID   string
	Name string
}

// Suggestion represents an autocomplete suggestion for JQL.
type Suggestion struct {
	Value       string
	DisplayName string
}

// SearchResult holds a page of search results.
type SearchResult struct {
	Issues   []Issue
	NextPage PageToken
	Total    int // -1 if unknown
}

// PageToken is an opaque cursor for pagination.
type PageToken string

// FieldMetadata describes a Jira field for JQL autocomplete.
type FieldMetadata struct {
	ID        string
	Name      string
	Operators []string
}
