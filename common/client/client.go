package client

import (
	"io"
	"net/http"
	"strings"

	"text/template"

	"github.com/google/uuid"
	"github.com/pauloo27/sonata/common/data"
)

type Response struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

type Client struct {
	httpClient http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.Client{},
	}
}

func (c *Client) Run(req *data.Request, params map[string]any) (*Response, error) {
	uriTemplate := req.URL

	var sb strings.Builder

	err := template.Must(
		template.New("url").
			Funcs(template.FuncMap{
				"randomUUID": randomUUID,
			}).
			Parse(uriTemplate),
	).
		Execute(&sb, params)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(string(req.Method), sb.String(), nil)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var body string

	if httpRes.Body != nil {
		bodyData, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}
		body = string(bodyData)
	}

	defer httpRes.Body.Close()

	return &Response{
		StatusCode: httpRes.StatusCode,
		Body:       body,
		Headers:    httpRes.Header,
	}, nil
}

func randomUUID() string {
	return uuid.New().String()
}
