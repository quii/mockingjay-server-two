package mockingjay

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type (
	Endpoints []Endpoint

	Endpoint struct {
		ID          uuid.UUID `json:"ID"`
		Description string    `json:"description,omitempty"`
		Request     Request   `json:"request"`
		Response    Response  `json:"response"`
		CDCs        []CDC
	}

	CDC struct {
		BaseURL   string `json:"baseURL"`
		Retries   int    `json:"retries"`
		TimeoutMS int    `json:"timeoutMS"`
	}

	Response struct {
		Status  int     `json:"status,omitempty"`
		Body    string  `json:"body,omitempty"`
		Headers Headers `json:"headers,omitempty"`
	}

	Request struct {
		Method    string  `json:"method,omitempty"`
		RegexPath string  `json:"regexPath,omitempty"`
		Path      string  `json:"path,omitempty"`
		Headers   Headers `json:"headers,omitempty"`
		Body      string  `json:"body,omitempty"`

		compiledRegex *regexp.Regexp
	}

	Headers map[string][]string
)

func (r *Request) MatchPath(path string) bool {
	if r.compiledRegex != nil {
		return r.compiledRegex.MatchString(path)
	}
	return r.Path == path
}

func (r *Request) Compile() error {
	if r.RegexPath != "" {
		rgx, err := regexp.Compile(r.RegexPath)
		if err != nil {
			return err
		}
		r.compiledRegex = rgx
	}
	return nil
}

func (r *Request) ToHTTPRequest(basePath string) *http.Request {
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

func (e Endpoints) Compile() error {
	for i := range e {
		if err := e[i].Request.Compile(); err != nil {
			return err
		}
		e[i].ID = uuid.New()
	}
	return nil
}
