package stub

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Request struct {
	Method    string  `json:"method,omitempty"`
	RegexPath string  `json:"regexPath,omitempty"`
	Path      string  `json:"path,omitempty"`
	Headers   Headers `json:"headers,omitempty"`
	Body      string  `json:"body,omitempty"`

	compiledRegex *regexp.Regexp
}

func NewRequestFromHTTP(req *http.Request) Request {
	var body string
	if req.Body != nil {
		defer req.Body.Close()
		bdy, _ := io.ReadAll(req.Body)
		body = string(bdy)
	}

	return Request{
		Method:  req.Method,
		Path:    req.URL.Path,
		Headers: Headers(req.Header),
		Body:    body,
	}
}

func (r *Request) MatchPath(path string) bool {
	if r.compiledRegex != nil {
		matchString := r.compiledRegex.MatchString(path)
		return matchString
	}
	return r.Path == path
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

func (r *Request) compile() error {
	if r.RegexPath != "" {
		rgx, err := regexp.Compile(r.RegexPath)
		if err != nil {
			return err
		}
		r.compiledRegex = rgx
	}
	return nil
}
