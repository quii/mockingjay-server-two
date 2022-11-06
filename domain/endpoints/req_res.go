package endpoints

import (
	"io"
	"net/http"
	"strings"
)

type Headers map[string][]string

type Response struct {
	Status  int     `json:"status,omitempty"`
	Body    string  `json:"body,omitempty"`
	Headers Headers `json:"headers,omitempty"`
}

type Request struct {
	Method  string  `json:"method,omitempty"`
	Path    string  `json:"path,omitempty"`
	Headers Headers `json:"headers,omitempty"`
	Body    string  `json:"body,omitempty"`
}

func (r Request) ToHTTPRequest(basePath string) *http.Request {
	req, _ := http.NewRequest(r.Method, basePath+r.Path, nil)

	if r.Body != "" {
		req.Body = io.NopCloser(strings.NewReader(r.Body))
	}

	for key, values := range r.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req
}

type Endpoint struct {
	Description string   `json:"description,omitempty"`
	Request     Request  `json:"request"`
	Response    Response `json:"response"`
}
