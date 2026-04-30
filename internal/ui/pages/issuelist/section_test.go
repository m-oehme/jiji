package issuelist

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/m-oehme/jiji/internal/config"
	"github.com/m-oehme/jiji/internal/jira"
)

func testSectionedIssues() []jira.Issue {
	customFieldsAndroid := map[string]any{
		"customfield_10001": "Android",
	}
	customFieldsIOS := map[string]any{
		"customfield_10001": "iOS",
	}

	return []jira.Issue{
		{Key: "TEST-1", Summary: "First issue", Priority: "High", Assignee: "Charlie", Status: "TODO", CustomFields: customFieldsAndroid},
		{Key: "TEST-2", Summary: "Second issue", Priority: "Medium", Assignee: "Charlie", Status: "TODO"},
		{Key: "TEST-3", Summary: "Third issue", Priority: "Medium", Assignee: "Charlie", Status: "TODO"},
		{Key: "TEST-4", Summary: "Forth issue", Priority: "Low", Assignee: "Charlie", Status: "InProgress"},
		{Key: "TEST-5", Summary: "Fifth issue", Priority: "Low", Assignee: "Charlie", Status: "Done", CustomFields: customFieldsAndroid},
		{Key: "TEST-6", Summary: "Sixth issue", Priority: "Low", Assignee: "Charlie", Status: "Review", CustomFields: customFieldsIOS},
	}
}

func TestBuildRows_DropUnmatched(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled:          true,
		Field:            "{{ .Status }}",
		Values:           []string{"TODO", "InProgress", "Done"},
		IncludeUnmatched: false,
	}

	want := []Row{
		{Header: "TODO"},
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
		{Header: "InProgress"},
		{Issue: &issues[3]},
		{Header: "Done"},
		{Issue: &issues[4]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_IncludeUnmatched(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled:          true,
		Field:            "{{ .Status }}",
		Values:           []string{"TODO", "InProgress", "Done"},
		IncludeUnmatched: true,
	}

	want := []Row{
		{Header: "TODO"},
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
		{Header: "InProgress"},
		{Issue: &issues[3]},
		{Header: "Done"},
		{Issue: &issues[4]},
		{Header: "Other"},
		{Issue: &issues[5]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_AutoDiscover(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled: true,
		Field:   "{{ .Status }}",
	}

	want := []Row{
		{Header: "Done"},
		{Issue: &issues[4]},
		{Header: "InProgress"},
		{Issue: &issues[3]},
		{Header: "Review"},
		{Issue: &issues[5]},
		{Header: "TODO"},
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_Disabled(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled: false,
	}

	want := []Row{
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
		{Issue: &issues[3]},
		{Issue: &issues[4]},
		{Issue: &issues[5]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_EmptyIssues(t *testing.T) {
	config := config.SectionsConfig{
		Enabled: false,
	}

	want := []Row{}

	got := BuildRows([]jira.Issue{}, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_SingleSection(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled: true,
		Field:   "{{ .Assignee }}",
	}

	want := []Row{
		{Header: "Charlie"},
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
		{Issue: &issues[3]},
		{Issue: &issues[4]},
		{Issue: &issues[5]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_CompositeTemplate(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled:          true,
		Field:            "{{ .Status }}-{{ .Priority }}",
		Values:           []string{"TODO-High", "TODO-Medium", "InProgress-Low"},
		IncludeUnmatched: false,
	}

	want := []Row{
		{Header: "TODO-High"},
		{Issue: &issues[0]},
		{Header: "TODO-Medium"},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
		{Header: "InProgress-Low"},
		{Issue: &issues[3]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_CustomFields(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled:          true,
		Field:            "{{ index .CustomFields \"customfield_10001\" }}",
		Values:           []string{"Android", "iOS"},
		IncludeUnmatched: false,
	}

	want := []Row{
		{Header: "Android"},
		{Issue: &issues[0]},
		{Issue: &issues[4]},
		{Header: "iOS"},
		{Issue: &issues[5]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_EmptySection(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled: true,
		Field:   "{{ .Status }}",
		Values:  []string{"TODO", "Backlog"},
	}

	want := []Row{
		{Header: "TODO"},
		{Issue: &issues[0]},
		{Issue: &issues[1]},
		{Issue: &issues[2]},
	}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}

func TestBuildRows_TemplateError(t *testing.T) {
	issues := testSectionedIssues()
	config := config.SectionsConfig{
		Enabled: true,
		Field:   "{{ ThisIsNotValid }}",
	}

	want := []Row{}

	got := BuildRows(issues, config)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sectioned isses mismatch (-want +got):\n%s", diff)
	}
}
