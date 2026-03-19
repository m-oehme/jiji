// Package jira defines the client interface and adapters for interacting with Jira.
package jira

import (
	"context"
	"encoding/json"
)

// Client defines the interface for interacting with a Jira instance.
// For v1, only the Cloud adapter implements this (ADR-004).
type Client interface {
	// Reading
	SearchIssues(ctx context.Context, jql string, fields []string, page PageToken) (*SearchResult, error)
	GetIssue(ctx context.Context, key string) (*Issue, error)
	GetComments(ctx context.Context, issueKey string) ([]Comment, error)
	GetFieldMetadata(ctx context.Context) ([]FieldMetadata, error)
	GetAutocompleteSuggestions(ctx context.Context, fieldName, fieldValue string) ([]Suggestion, error)

	// Editing
	GetTransitions(ctx context.Context, key string) ([]Transition, error)
	TransitionIssue(ctx context.Context, key string, transitionID string) error
	UpdateIssue(ctx context.Context, key string, fields map[string]any) error
	AddComment(ctx context.Context, key string, adfBody json.RawMessage) error
	UpdateComment(ctx context.Context, key, commentID string, adfBody json.RawMessage) error
	GetLabels(ctx context.Context) ([]string, error)
}
