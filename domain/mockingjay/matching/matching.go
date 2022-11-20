package matching

import (
	"net/http"
	"net/textproto"

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
	got := mockingjay.NewRequestFromHTTP(req)

	return func(e mockingjay.Endpoint) RequestMatch {
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

func matchBody(a, b mockingjay.Request) bool {
	return a.Body == b.Body
}

func matchPath(a, b mockingjay.Request) bool {
	return a.MatchPath(b.Path)
}

func matchMethod(a, b mockingjay.Request) bool {
	return a.Method == b.Method
}

func matchHeaders(a, b mockingjay.Request) bool {
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
