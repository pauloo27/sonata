package client

import (
	"html"
	"net/url"
	"os"
	"text/template"

	"github.com/google/uuid"
)

var (
	GetEnv = os.Getenv
)

var (
	templateFuncs = template.FuncMap{
		"randomUUID":  randomUUID,
		"queryEscape": queryEscape,
		"pathEscape":  pathEscape,
		"htmlEscape":  htmlEscape,
		"env":         env,
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

func env(s string) string {
	return GetEnv(s)
}

func orDefault(s string, defaultValue string) string {
	if s == "" {
		return defaultValue
	}
	return s
}
