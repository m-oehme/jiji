package issuelist

import (
	"bytes"
	"maps"
	"slices"
	"sort"
	"text/template"

	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
)

type Row struct {
	Header string
	Issue  *jira.Issue
}

func BuildRows(issues []jira.Issue, cfg config.SectionsConfig) []Row {
	if !cfg.Enabled {
		rows := make([]Row, len(issues))
		for idx := range issues {
			rows[idx] = Row{Issue: &issues[idx]}
		}
		return rows
	}

	template, err := template.New("section_field").Parse(cfg.Field)
	if err != nil {
		return make([]Row, 0)
	}

	var buffer bytes.Buffer
	resolvedIssues := make(map[string][]Row)

	for idx, issue := range issues {
		err = template.Execute(&buffer, issue)
		if err != nil {
			buffer.WriteString("")
		}

		key := buffer.String()
		row := Row{Issue: &issues[idx]}
		resolvedIssues[key] = append(resolvedIssues[key], row)

		buffer.Reset()
	}

	var sectionOrder []string

	if len(cfg.Values) != 0 {
		sectionOrder = cfg.Values
	} else {
		keys := slices.Collect(maps.Keys(resolvedIssues))
		sort.Strings(keys)
		sectionOrder = keys
	}

	var rows []Row
	for _, section := range sectionOrder {
		issueSlice := resolvedIssues[section]
		if len(issueSlice) != 0 {
			header := Row{Header: section}
			rows = append(rows, header)
			rows = append(rows, issueSlice...)
		}
	}

	if cfg.IncludeUnmatched {
		rows = append(rows, Row{Header: "Other"})
		for key, sectionIssue := range resolvedIssues {
			if !slices.Contains(sectionOrder, key) {
				rows = append(rows, sectionIssue...)
			}
		}
	}

	return rows
}
