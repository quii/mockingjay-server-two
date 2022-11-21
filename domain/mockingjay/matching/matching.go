package matching

import (
	"net/http"
	"net/textproto"

	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

type RequestMatch struct {
	Endpoint http2.Endpoint
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

func newMatcher(req *http.Request) func(http2.Endpoint) RequestMatch {
	got := http2.NewRequestFromHTTP(req)

	return func(e http2.Endpoint) RequestMatch {
		match := RequestMatch{
			Endpoint: e,
			Match: Match{
				Path:    matchPath(e.Request, got),
				Method:  matchMethod(e.Request, got),
				Headers: matchHeaders(e.Request, got),
				Body:    matchBody(e.Request, got),
			},
		}
		return match
	}
}

func matchBody(a, b http2.Request) bool {
	return a.Body == b.Body
}

func matchPath(a, b http2.Request) bool {
	return a.MatchPath(b.Path)
}

func matchMethod(a, b http2.Request) bool {
	return a.Method == b.Method
}

func matchHeaders(a, b http2.Request) bool {
	headersMatch := len(a.Headers) == 0

	for key, values := range a.Headers {
		for _, valuesInIncomingRequestHeader := range b.Headers[textproto.CanonicalMIMEHeaderKey(key)] {
			for _, valuesInEndpoint := range values {
				if valuesInIncomingRequestHeader == valuesInEndpoint {
					headersMatch = true
				}
			}
		}
	}
	return headersMatch
}
