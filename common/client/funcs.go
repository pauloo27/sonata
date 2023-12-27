package client

import (
	"text/template"

	"github.com/google/uuid"
)

var (
	templateFuncs = template.FuncMap{
		"randomUUID": randomUUID,
	}
)

func randomUUID() string {
	return uuid.New().String()
}
