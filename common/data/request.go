package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type BodyType string

const (
	BodyTypeNone BodyType = "none"
	BodyTypeJSON BodyType = "json"
	BodyTypeForm BodyType = "form"
	BodyTypeText BodyType = "text"
)

var (
	BodyTypeExtensions = map[BodyType]string{
		BodyTypeNone: "",
		BodyTypeJSON: "json",
		BodyTypeForm: "form",
		BodyTypeText: "txt",
	}

	BodyTypes = [...]BodyType{
		BodyTypeNone,
		BodyTypeJSON,
		BodyTypeForm,
		BodyTypeText,
	}
)

type HTTPMethod string

const (
	HTTPMethodGet    HTTPMethod = "GET"
	HTTPMethodPost   HTTPMethod = "POST"
	HTTPMethodPut    HTTPMethod = "PUT"
	HTTPMethodDelete HTTPMethod = "DELETE"
	HTTPMethodPatch  HTTPMethod = "PATCH"
	HTTPMethodHead   HTTPMethod = "HEAD"
	HTTPMethodOption HTTPMethod = "OPTION"
)

var (
	HTTPMethods = [...]HTTPMethod{
		HTTPMethodGet,
		HTTPMethodPost,
		HTTPMethodPut,
		HTTPMethodDelete,
		HTTPMethodPatch,
		HTTPMethodHead,
		HTTPMethodOption,
	}
)

type Request struct {
	Name        string     `json:"-"`
	Description string     `json:"description"`
	Method      HTTPMethod `json:"method"`

	URL      string            `json:"url"`
	BodyType BodyType          `json:"body_type"`
	Body     string            `json:"body_template"`
	Headers  map[string]string `json:"headers"`

	Version string `json:"version"`

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

		path: fmt.Sprintf("%s/%s.request.json", p.rootDir, name),
		p:    p,
	}
}

func (r *Request) Save() error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return os.WriteFile(r.path, data, 420)
}

func (r *Request) Rename(newName string) error {
	if err := os.Remove(r.path); err != nil {
		return err
	}

	r.Name = newName
	r.path = fmt.Sprintf("%s/%s.request.json", r.p.rootDir, newName)

	return r.Save()
}

func (r *Request) Delete() error {
	return os.Remove(r.path)
}

func (r *Request) Clone() *Request {
	return &Request{
		Name:        r.Name,
		Description: r.Description,
		Method:      r.Method,
		URL:         r.URL,
		BodyType:    r.BodyType,
		Body:        r.Body,
		Headers:     r.Headers,
		Version:     r.Version,

		path: r.path,
		p:    r.p,
	}
}
