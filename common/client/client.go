package client

import (
	"io"
	"net/http"
	"strings"

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

func (c *Client) Run(req *data.Request, variables map[string]string) (*Response, error) {
	uri, err := ExecuteTemplate(
		"req-uri-tmpl",
		req.URL,
		variables,
	)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader

	if req.Body != "" {
		body, err := ExecuteTemplate(
			"req-body-tmpl",
			req.Body,
			variables,
		)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(body)
	}

	httpReq, err := http.NewRequest(string(req.Method), uri, bodyReader)
	if err != nil {
		return nil, err
	}

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var resBody string

	if httpRes.Body != nil {
		resBodyData, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, err
		}
		resBody = string(resBodyData)
	}

	defer httpRes.Body.Close()

	return &Response{
		StatusCode: httpRes.StatusCode,
		Body:       resBody,
		Headers:    httpRes.Header,
	}, nil
}
