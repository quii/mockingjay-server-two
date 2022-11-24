package matching

import (
	"net/http"
	"net/textproto"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type RequestMatch struct {
	Endpoint stub.Endpoint
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

func newMatcher(req *http.Request) func(stub.Endpoint) RequestMatch {
	got := stub.NewRequestFromHTTP(req)

	return func(e stub.Endpoint) RequestMatch {
		match := RequestMatch{
			Endpoint: e,
			Match: Match{
				Path:    matchPath(e.Request, got),
				Method:  matchMethod(e.Request, got),
				Headers: MatchHeaders(e.Request.Headers, got.Headers),
				Body:    matchBody(e.Request, got),
			},
		}
		return match
	}
}

func matchBody(a, b stub.Request) bool {
	return a.Body == b.Body
}

func matchPath(a, b stub.Request) bool {
	return a.MatchPath(b.Path)
}

func matchMethod(a, b stub.Request) bool {
	return a.Method == b.Method
}

func MatchHeaders(a, b stub.Headers) bool {
	headersMatch := len(a) == 0

	for key, values := range a {
		for _, valuesInIncomingRequestHeader := range b[textproto.CanonicalMIMEHeaderKey(key)] {
			for _, valuesInEndpoint := range values {
				if valuesInIncomingRequestHeader == valuesInEndpoint {
					headersMatch = true
				}
			}
		}
	}
	return headersMatch
}
