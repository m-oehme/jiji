package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	v3 "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// CloudAdapter implements Client for Jira Cloud using the v3 REST API.
type CloudAdapter struct {
	client *v3.Client
	host   string
}

// NewCloudAdapter creates a new Jira Cloud client.
func NewCloudAdapter(host, email, token string) (*CloudAdapter, error) {
	client, err := v3.New(nil, host)
	if err != nil {
		return nil, fmt.Errorf("creating Jira client: %w", err)
	}
	client.Auth.SetBasicAuth(email, token)
	return &CloudAdapter{client: client, host: host}, nil
}

func (a *CloudAdapter) SearchIssues(ctx context.Context, jql string, fields []string, page PageToken) (*SearchResult, error) {
	result, resp, err := a.client.Issue.Search.SearchJQL(ctx, jql, fields, nil, 50, string(page))
	if err != nil {
		return nil, wrapError(resp, err)
	}

	issues := make([]Issue, 0, len(result.Issues))
	for _, iss := range result.Issues {
		issues = append(issues, mapIssue(iss))
	}

	return &SearchResult{
		Issues:   issues,
		NextPage: PageToken(result.NextPageToken),
		Total:    result.Total,
	}, nil
}

func (a *CloudAdapter) GetIssue(ctx context.Context, key string) (*Issue, error) {
	iss, resp, err := a.client.Issue.Get(ctx, key, nil, nil)
	if err != nil {
		return nil, wrapError(resp, err)
	}
	issue := mapIssue(iss)
	return &issue, nil
}

func (a *CloudAdapter) GetComments(ctx context.Context, issueKey string) ([]Comment, error) {
	page, resp, err := a.client.Issue.Comment.Gets(ctx, issueKey, "created", nil, 0, 100)
	if err != nil {
		return nil, wrapError(resp, err)
	}

	comments := make([]Comment, 0, len(page.Comments))
	for _, c := range page.Comments {
		comments = append(comments, mapComment(c))
	}
	return comments, nil
}

func (a *CloudAdapter) GetFieldMetadata(ctx context.Context) ([]FieldMetadata, error) {
	// go-atlassian doesn't wrap /jql/autocompletedata, so we use a raw request.
	req, err := a.client.NewRequest(ctx, http.MethodGet, "rest/api/3/jql/autocompletedata", "", nil)
	if err != nil {
		return nil, fmt.Errorf("creating autocomplete request: %w", err)
	}

	var result struct {
		VisibleFieldNames []struct {
			Value       string   `json:"value"`
			DisplayName string   `json:"displayName"`
			Operators   []string `json:"operators"`
		} `json:"visibleFieldNames"`
	}

	resp, err := a.client.Call(req, &result)
	if err != nil {
		return nil, wrapError(resp, err)
	}

	fields := make([]FieldMetadata, 0, len(result.VisibleFieldNames))
	for _, f := range result.VisibleFieldNames {
		fields = append(fields, FieldMetadata{
			ID:        f.Value,
			Name:      f.DisplayName,
			Operators: f.Operators,
		})
	}
	return fields, nil
}

func (a *CloudAdapter) GetAutocompleteSuggestions(ctx context.Context, fieldName, fieldValue string) ([]Suggestion, error) {
	endpoint := fmt.Sprintf("rest/api/3/jql/autocompletedata/suggestions?fieldName=%s&fieldValue=%s", fieldName, fieldValue)
	req, err := a.client.NewRequest(ctx, http.MethodGet, endpoint, "", nil)
	if err != nil {
		return nil, fmt.Errorf("creating suggestions request: %w", err)
	}

	var result struct {
		Results []struct {
			Value       string `json:"value"`
			DisplayName string `json:"displayName"`
		} `json:"results"`
	}

	resp, err := a.client.Call(req, &result)
	if err != nil {
		return nil, wrapError(resp, err)
	}

	suggestions := make([]Suggestion, 0, len(result.Results))
	for _, s := range result.Results {
		suggestions = append(suggestions, Suggestion{
			Value:       s.Value,
			DisplayName: s.DisplayName,
		})
	}
	return suggestions, nil
}

func (a *CloudAdapter) GetTransitions(ctx context.Context, key string) ([]Transition, error) {
	result, resp, err := a.client.Issue.Transitions(ctx, key)
	if err != nil {
		return nil, wrapError(resp, err)
	}

	transitions := make([]Transition, 0, len(result.Transitions))
	for _, t := range result.Transitions {
		transitions = append(transitions, Transition{
			ID:   t.ID,
			Name: t.Name,
		})
	}
	return transitions, nil
}

func (a *CloudAdapter) TransitionIssue(ctx context.Context, key string, transitionID string) error {
	resp, err := a.client.Issue.Move(ctx, key, transitionID, nil)
	if err != nil {
		return wrapError(resp, err)
	}
	return nil
}

func (a *CloudAdapter) UpdateIssue(ctx context.Context, key string, fields map[string]any) error {
	payload := &models.IssueScheme{
		Fields: &models.IssueFieldsScheme{},
	}

	// Map known fields to the typed struct; remaining go to custom fields.
	custom := &models.CustomFields{}
	hasCustom := false
	for k, v := range fields {
		switch k {
		case "summary":
			if s, ok := v.(string); ok {
				payload.Fields.Summary = s
			}
		case "labels":
			if labels, ok := v.([]string); ok {
				payload.Fields.Labels = labels
			}
		default:
			// Treat as custom field
			if raw, err := json.Marshal(v); err == nil {
				custom.Fields = append(custom.Fields, map[string]interface{}{k: json.RawMessage(raw)})
				hasCustom = true
			}
		}
	}

	var customPtr *models.CustomFields
	if hasCustom {
		customPtr = custom
	}

	resp, err := a.client.Issue.Update(ctx, key, true, payload, customPtr, nil)
	if err != nil {
		return wrapError(resp, err)
	}
	return nil
}

func (a *CloudAdapter) AddComment(ctx context.Context, key string, adfBody json.RawMessage) error {
	var body models.CommentNodeScheme
	if err := json.Unmarshal(adfBody, &body); err != nil {
		return fmt.Errorf("invalid ADF body: %w", err)
	}

	payload := &models.CommentPayloadScheme{Body: &body}
	_, resp, err := a.client.Issue.Comment.Add(ctx, key, payload, nil)
	if err != nil {
		return wrapError(resp, err)
	}
	return nil
}

func (a *CloudAdapter) UpdateComment(ctx context.Context, key, commentID string, adfBody json.RawMessage) error {
	var body models.CommentNodeScheme
	if err := json.Unmarshal(adfBody, &body); err != nil {
		return fmt.Errorf("invalid ADF body: %w", err)
	}

	payload := &models.CommentPayloadScheme{Body: &body}
	_, resp, err := a.client.Issue.Comment.Update(ctx, key, commentID, payload, nil)
	if err != nil {
		return wrapError(resp, err)
	}
	return nil
}

func (a *CloudAdapter) GetLabels(ctx context.Context) ([]string, error) {
	var allLabels []string
	startAt := 0
	for {
		page, resp, err := a.client.Issue.Label.Gets(ctx, startAt, 1000)
		if err != nil {
			return nil, wrapError(resp, err)
		}
		allLabels = append(allLabels, page.Values...)
		if page.IsLast || len(page.Values) == 0 {
			break
		}
		startAt += len(page.Values)
	}
	return allLabels, nil
}

// mapIssue converts a go-atlassian IssueScheme to our domain Issue.
func mapIssue(iss *models.IssueScheme) Issue {
	issue := Issue{
		Key: iss.Key,
	}

	f := iss.Fields
	if f == nil {
		return issue
	}

	issue.Summary = f.Summary

	// Marshal ADF description back to json.RawMessage
	if f.Description != nil {
		if raw, err := json.Marshal(f.Description); err == nil {
			issue.Description = raw
		}
	}

	if f.Status != nil {
		issue.Status = f.Status.Name
		issue.StatusID = f.Status.ID
	}
	if f.Priority != nil {
		issue.Priority = f.Priority.Name
	}
	if f.IssueType != nil {
		issue.Type = f.IssueType.Name
	}
	if f.Assignee != nil {
		issue.Assignee = f.Assignee.DisplayName
	}
	if f.Reporter != nil {
		issue.Reporter = f.Reporter.DisplayName
	}

	issue.Labels = f.Labels

	for _, comp := range f.Components {
		if comp != nil {
			issue.Components = append(issue.Components, comp.Name)
		}
	}

	for _, v := range f.FixVersions {
		if v != nil {
			issue.FixVersions = append(issue.FixVersions, v.Name)
		}
	}

	if f.Created != nil {
		issue.Created = time.Time(*f.Created)
	}
	if f.Updated != nil {
		issue.Updated = time.Time(*f.Updated)
	}

	// Subtasks
	for _, sub := range f.Subtasks {
		if sub != nil {
			link := IssueLink{Key: sub.Key}
			if sub.Fields != nil {
				link.Summary = sub.Fields.Summary
				if sub.Fields.Status != nil {
					link.Status = sub.Fields.Status.Name
				}
			}
			link.LinkType = "subtask"
			issue.Subtasks = append(issue.Subtasks, link)
		}
	}

	// Issue links
	for _, l := range f.IssueLinks {
		if l == nil {
			continue
		}
		var linked *models.LinkedIssueScheme
		var linkType string
		if l.OutwardIssue != nil {
			linked = l.OutwardIssue
			if l.Type != nil {
				linkType = l.Type.Outward
			}
		} else if l.InwardIssue != nil {
			linked = l.InwardIssue
			if l.Type != nil {
				linkType = l.Type.Inward
			}
		}
		if linked != nil {
			link := IssueLink{
				Key:      linked.Key,
				LinkType: linkType,
			}
			if linked.Fields != nil {
				link.Summary = linked.Fields.Summary
				if linked.Fields.Status != nil {
					link.Status = linked.Fields.Status.Name
				}
			}
			issue.Links = append(issue.Links, link)
		}
	}

	return issue
}

// mapComment converts a go-atlassian comment to our domain Comment.
func mapComment(c *models.IssueCommentScheme) Comment {
	comment := Comment{
		ID: c.ID,
	}
	if c.Author != nil {
		comment.Author = c.Author.DisplayName
	}
	if c.Body != nil {
		if raw, err := json.Marshal(c.Body); err == nil {
			comment.Body = raw
		}
	}
	comment.Created = parseJiraTime(c.Created)
	comment.Updated = parseJiraTime(c.Updated)
	return comment
}

// parseJiraTime parses Jira's datetime format.
func parseJiraTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	// Jira uses ISO 8601 with timezone offset
	formats := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05.000Z0700",
		time.RFC3339,
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// wrapError converts HTTP response codes into user-friendly error messages.
func wrapError(resp *models.ResponseScheme, err error) error {
	if resp == nil {
		if err != nil && (strings.Contains(err.Error(), "dial") || strings.Contains(err.Error(), "connection")) {
			return fmt.Errorf("connection failed: %w", err)
		}
		return err
	}

	switch resp.Code {
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed, check --email and --token: %w", err)
	case http.StatusForbidden:
		return fmt.Errorf("permission denied: %w", err)
	case http.StatusNotFound:
		return fmt.Errorf("not found: %w", err)
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		if secs, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
			return &RateLimitError{RetryAfter: time.Duration(secs) * time.Second, Err: err}
		}
		return fmt.Errorf("rate limited: %w", err)
	default:
		return err
	}
}

// RateLimitError indicates the API returned 429 with a Retry-After duration.
type RateLimitError struct {
	RetryAfter time.Duration
	Err        error
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited, retry after %s: %s", e.RetryAfter, e.Err)
}

func (e *RateLimitError) Unwrap() error {
	return e.Err
}

// Compile-time check that CloudAdapter implements Client.
var _ Client = (*CloudAdapter)(nil)
