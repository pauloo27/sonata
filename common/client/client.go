package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pauloo27/sonata/common/data"
)

type Response struct {
	StatusCode int
	Body       string
	Headers    http.Header

	CalledURL string
	Time      time.Duration
}

type Client struct {
	httpClient http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.Client{},
	}
}

func UseMap(variables map[string]string) {
	GetEnv = func(name string) string {
		return variables[name]
	}
}

func UseEnvFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	variables, err := godotenv.Parse(file)
	if err != nil {
		return err
	}

	GetEnv = func(name string) string {
		return variables[name]
	}

	return nil
}

var (
	bodyReaderFnMap = map[data.BodyType]func(*data.Request, map[string]string) (io.Reader, error){
		data.BodyTypeJSON: newTextBodyReader,
		data.BodyTypeText: newTextBodyReader,
		data.BodyTypeForm: newFormBodyReader,
	}
)

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
	readerFn, found := bodyReaderFnMap[req.BodyType]

	if found {
		bodyReader, err = readerFn(req, variables)
		if err != nil {
			return nil, err
		}
	}

	httpReq, err := http.NewRequest(string(req.Method), uri, bodyReader)
	if err != nil {
		return nil, err
	}

	switch req.BodyType {
	case data.BodyTypeJSON:
		httpReq.Header.Add("Content-Type", "application/json")
	case data.BodyTypeForm:
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	for key, valueTmpl := range req.Headers {
		finalValue, err := ExecuteTemplate(
			"req-header-tmpl",
			valueTmpl,
			variables,
		)
		if err != nil {
			return nil, err
		}
		httpReq.Header.Add(key, finalValue)
	}

	start := time.Now()

	httpRes, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	took := time.Since(start)

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
		CalledURL:  uri,
		StatusCode: httpRes.StatusCode,
		Body:       resBody,
		Headers:    httpRes.Header,
		Time:       took,
	}, nil
}

func newFormBodyReader(req *data.Request, variables map[string]string) (io.Reader, error) {
	if req.Body == "" {
		return nil, nil
	}

	var rawValues map[string]string

	err := json.Unmarshal([]byte(req.Body), &rawValues)
	if err != nil {
		return nil, err
	}

	values := url.Values{}

	for key, valueTmpl := range rawValues {
		value, err := ExecuteTemplate(
			"req-body-tmpl",
			valueTmpl,
			variables,
		)
		if err != nil {
			return nil, err
		}
		values.Add(key, value)
	}

	return strings.NewReader(values.Encode()), nil
}

func newTextBodyReader(req *data.Request, variables map[string]string) (io.Reader, error) {
	if req.Body == "" {
		return nil, nil
	}

	body, err := ExecuteTemplate(
		"req-body-tmpl",
		req.Body,
		variables,
	)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(body), nil
}
