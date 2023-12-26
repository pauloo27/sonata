package client

import (
	"io"
	"net/http"

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

func (c *Client) Run(req *data.Request) (*Response, error) {
	// TODO: body and headers
	httpReq, err := http.NewRequest(string(req.Method), req.URL, nil)
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
