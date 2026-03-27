package jira

import (
	"bytes"
	"text/template"
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
