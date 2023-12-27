package client

import (
	"strings"
	"text/template"
)

func ExecuteTemplate(key string, tmplStr string, variables map[string]string) (string, error) {
	tmpl, err := LoadTemplate(key, tmplStr)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, variables); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func LoadTemplate(key string, tmplStr string) (*template.Template, error) {
	tmpl, err := template.New(key).
		Funcs(templateFuncs).
		Parse(tmplStr)

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
