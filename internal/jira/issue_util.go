package jira

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/m-oehme/jiji/internal/config"
)

func (m *Issue) Format(format string) (string, error) {
	template, err := template.New("issue").Parse(format)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = template.Execute(&buffer, m)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (m *Issue) IssueURL(config config.Config) string {
	host := strings.Trim(config.Jira.Host, "/")
	return fmt.Sprintf("%s/browse/%s", host, m.Key)
}
