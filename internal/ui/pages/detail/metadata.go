// Package detail renders the right pane: metadata banner + scrollable content (ADR-005).
package detail

import (
	"fmt"
	"strings"

	"github.com/m-oehme/jiji/internal/jira"
	"github.com/m-oehme/jiji/internal/ui/styles"
	lipgloss "charm.land/lipgloss/v2"
)

// metadataBannerHeight is the fixed height of the metadata section.
const metadataBannerHeight = 3

// renderMetadata renders the compact metadata banner for an issue.
func renderMetadata(issue *jira.Issue, width int, s *styles.Styles) string {
	if issue == nil {
		return strings.Repeat(" ", width)
	}

	pairs := []struct{ key, val string }{
		{"Status", issue.Status},
		{"Priority", issue.Priority},
		{"Assignee", issue.Assignee},
		{"Type", issue.Type},
	}

	if issue.Sprint != "" {
		pairs = append(pairs, struct{ key, val string }{"Sprint", issue.Sprint})
	}

	var parts []string
	for _, p := range pairs {
		if p.val == "" {
			continue
		}
		kv := fmt.Sprintf("%s: %s",
			s.MetadataKey.Render(p.key),
			s.MetadataValue.Render(p.val),
		)
		parts = append(parts, kv)
	}

	// Layout: join with separator, wrap if needed
	sep := s.Dimmed.Render(" │ ")
	line := strings.Join(parts, sep)

	// Labels on second line if present
	var labelLine string
	if len(issue.Labels) > 0 {
		labelLine = s.MetadataKey.Render("Labels") + ": " +
			s.MetadataValue.Render(strings.Join(issue.Labels, ", "))
	}

	content := line
	if labelLine != "" {
		content += "\n" + labelLine
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(metadataBannerHeight).
		PaddingLeft(1).
		Render(content)
}
