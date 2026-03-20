package detail

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/styles"
)

// buildContent assembles the scrollable content string for an issue.
// Phase 2: renders description as plain text (ADF rendering comes Phase 4).
func buildContent(issue *jira.Issue, comments []jira.Comment, s *styles.Styles) string {
	if issue == nil {
		return s.Dimmed.Render("No issue selected")
	}

	var sections []string

	// Summary heading
	sections = append(sections, s.Heading.Render(issue.Summary))
	sections = append(sections, "")

	// Description
	desc := extractADFText(issue.Description)
	if desc == "" {
		desc = s.Dimmed.Render("(no description)")
	}
	sections = append(sections, desc)
	sections = append(sections, "")

	// Subtasks
	if len(issue.Subtasks) > 0 {
		sections = append(sections, s.Heading.Render("── Subtasks ──"))
		for _, sub := range issue.Subtasks {
			status := s.Dimmed.Render("[" + sub.Status + "]")
			sections = append(sections, fmt.Sprintf("  %s %s %s", sub.Key, sub.Summary, status))
		}
		sections = append(sections, "")
	}

	// Linked Issues
	if len(issue.Links) > 0 {
		sections = append(sections, s.Heading.Render("── Linked Issues ──"))
		for _, link := range issue.Links {
			rel := s.Dimmed.Render(link.LinkType)
			status := s.Dimmed.Render("[" + link.Status + "]")
			sections = append(sections, fmt.Sprintf("  %s %s %s %s", rel, link.Key, link.Summary, status))
		}
		sections = append(sections, "")
	}

	// Comments
	if len(comments) > 0 {
		sections = append(sections, s.Heading.Render("── Comments ──"))
		for _, c := range comments {
			header := fmt.Sprintf("  %s  %s",
				s.MetadataKey.Render(c.Author),
				s.Dimmed.Render(c.Created.Format("2006-01-02 15:04")),
			)
			body := extractADFText(c.Body)
			if body == "" {
				body = s.Dimmed.Render("(empty)")
			}
			sections = append(sections, header)
			sections = append(sections, "  "+body)
			sections = append(sections, "")
		}
	}

	return strings.Join(sections, "\n")
}

// extractADFText does a naive extraction of text from ADF JSON.
// Phase 2: just extracts text node values. Full ADF rendering comes Phase 4.
func extractADFText(adf json.RawMessage) string {
	if len(adf) == 0 {
		return ""
	}

	var doc struct {
		Content []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"content"`
	}

	if err := json.Unmarshal(adf, &doc); err != nil {
		// Fall back to raw string for non-ADF content
		var raw string
		if json.Unmarshal(adf, &raw) == nil {
			return raw
		}
		return string(adf)
	}

	var lines []string
	for _, block := range doc.Content {
		var texts []string
		for _, inline := range block.Content {
			if inline.Text != "" {
				texts = append(texts, inline.Text)
			}
		}
		if len(texts) > 0 {
			lines = append(lines, strings.Join(texts, ""))
		}
	}

	return strings.Join(lines, "\n")
}
