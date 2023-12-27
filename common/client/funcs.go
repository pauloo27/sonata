package client

import (
	"html"
	"net/url"
	"text/template"

	"github.com/google/uuid"
)

var (
	templateFuncs = template.FuncMap{
		"randomUUID":  randomUUID,
		"queryEscape": queryEscape,
		"pathEscape":  pathEscape,
		"htmlEscape":  htmlEscape,
	}
)

func randomUUID() string {
	return uuid.New().String()
}

func queryEscape(s string) string {
	return url.QueryEscape(s)
}

func pathEscape(s string) string {
	return url.PathEscape(s)
}

func htmlEscape(s string) string {
	return html.EscapeString(s)
}
