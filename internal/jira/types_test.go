package jira

import (
	"encoding/json"
	"testing"
)

func TestIssue_RawADFDescription(t *testing.T) {
	// Simulate unmarshaling a Jira API response with ADF description
	adf := json.RawMessage(`{
		"version": 1,
		"type": "doc",
		"content": [
			{
				"type": "paragraph",
				"content": [
					{"type": "text", "text": "Hello world"}
				]
			}
		]
	}`)

	issue := Issue{
		Key:         "TEST-123",
		Summary:     "Test issue",
		Description: adf,
		Status:      "In Progress",
		Priority:    "High",
	}

	// Verify raw ADF is preserved
	var doc map[string]any
	if err := json.Unmarshal(issue.Description, &doc); err != nil {
		t.Fatalf("failed to unmarshal ADF: %v", err)
	}
	if doc["type"] != "doc" {
		t.Errorf("expected ADF type 'doc', got %v", doc["type"])
	}
	if doc["version"] != float64(1) {
		t.Errorf("expected ADF version 1, got %v", doc["version"])
	}
}

func TestComment_RawADFBody(t *testing.T) {
	adf := json.RawMessage(`{
		"version": 1,
		"type": "doc",
		"content": [
			{
				"type": "paragraph",
				"content": [
					{"type": "text", "text": "A comment"}
				]
			}
		]
	}`)

	comment := Comment{
		ID:     "12345",
		Author: "John Doe",
		Body:   adf,
	}

	var doc map[string]any
	if err := json.Unmarshal(comment.Body, &doc); err != nil {
		t.Fatalf("failed to unmarshal comment ADF: %v", err)
	}
	content := doc["content"].([]any)
	para := content[0].(map[string]any)
	if para["type"] != "paragraph" {
		t.Errorf("expected paragraph, got %v", para["type"])
	}
}

func TestSearchResult_Total(t *testing.T) {
	sr := SearchResult{
		Issues:   []Issue{{Key: "TEST-1"}, {Key: "TEST-2"}},
		NextPage: PageToken("abc123"),
		Total:    42,
	}

	if len(sr.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(sr.Issues))
	}
	if sr.NextPage != "abc123" {
		t.Errorf("expected page token 'abc123', got %q", sr.NextPage)
	}
	if sr.Total != 42 {
		t.Errorf("expected total 42, got %d", sr.Total)
	}
}

func TestIssue_JSONMarshal(t *testing.T) {
	adf := json.RawMessage(`{"type":"doc","version":1,"content":[]}`)
	issue := Issue{
		Key:         "PROJ-456",
		Summary:     "JSON test",
		Description: adf,
		Labels:      []string{"bug", "urgent"},
		Components:  []string{"backend"},
	}

	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded Issue
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Key != issue.Key {
		t.Errorf("key mismatch: %q != %q", decoded.Key, issue.Key)
	}

	// Verify ADF roundtrips through json.RawMessage
	var origDoc, decodedDoc map[string]any
	if err := json.Unmarshal(issue.Description, &origDoc); err != nil {
		t.Fatalf("unmarshal original ADF: %v", err)
	}
	if err := json.Unmarshal(decoded.Description, &decodedDoc); err != nil {
		t.Fatalf("unmarshal decoded ADF: %v", err)
	}

	if origDoc["type"] != decodedDoc["type"] {
		t.Error("ADF type changed after roundtrip")
	}
}
