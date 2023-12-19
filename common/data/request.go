package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type BodyType string
type HTTPMethod string

const (
	BodyTypeJSON BodyType = "json"
	BodyTypeForm BodyType = "form"
	BodyTypeText BodyType = "text"
)

type Request struct {
	Name        string     `json:"-"`
	Description string     `json:"description"`
	Method      HTTPMethod `json:"method"`

	URL      string            `json:"url"`
	BodyType BodyType          `json:"body_type"`
	Body     string            `json:"body_template"`
	Version  string            `json:"version"`
	Headers  map[string]string `json:"headers"`

	path string   `json:"-"`
	p    *Project `json:"-"`
}

func (p *Project) NewRequest(
	name string, description string, method HTTPMethod,
	url string, bodyType BodyType, bodyTemplate string,
) *Request {
	return &Request{
		Name:        name,
		Description: description,
		Method:      method,
		URL:         url,
		BodyType:    bodyType,
		Body:        bodyTemplate,
		Version:     CurrentVersion,

		path: fmt.Sprintf("%s/%s.json", p.rootDir, name),
		p:    p,
	}
}

func (r *Request) Save() error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.json", r.p.rootDir, r.Name), data, 420)
}
