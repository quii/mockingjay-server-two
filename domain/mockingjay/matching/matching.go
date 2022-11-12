package matching

import (
	"io"
	"net/http"
	"net/textproto"
	"regexp"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

type RequestMatch struct {
	Endpoint mockingjay.Endpoint
	Match    Match `json:"match"`
}

func (r RequestMatch) Matched() bool {
	return r.Match.Path && r.Match.Method && r.Match.Headers && r.Match.Body
}

type Match struct {
	Path    bool `json:"path"`
	Method  bool `json:"method"`
	Headers bool `json:"headers"`
	Body    bool `json:"body"`
}

func newMatcher(req *http.Request) func(mockingjay.Endpoint) RequestMatch {
	var body []byte
	if req.Body != nil {
		defer req.Body.Close()
		body, _ = io.ReadAll(req.Body)
	}

	return func(e mockingjay.Endpoint) RequestMatch {
		return RequestMatch{
			Endpoint: e,
			Match: Match{
				Path:    matchPath(e, req),
				Method:  e.Request.Method == req.Method,
				Headers: matchHeaders(e, req.Header),
				Body:    string(body) == e.Request.Body,
			},
		}
	}
}

func matchPath(e mockingjay.Endpoint, req *http.Request) bool {
	if e.Request.RegexPath != "" {
		reg, err := regexp.Compile(e.Request.RegexPath)
		if err != nil {
			panic(err) //todo: move all this compilation stuff to loading
		}
		return reg.MatchString(req.URL.Path)
	}
	return e.Request.Path == req.URL.String()
}

func matchHeaders(e mockingjay.Endpoint, incomingHeaders http.Header) bool {
	headersMatch := len(e.Request.Headers) == 0

	for key, values := range e.Request.Headers {
		for _, valuesInIncomingRequestHeader := range incomingHeaders[textproto.CanonicalMIMEHeaderKey(key)] {
			for _, valuesInEndpoint := range values {
				if valuesInIncomingRequestHeader == valuesInEndpoint {
					headersMatch = true
				}
			}
		}
	}
	return headersMatch
}
