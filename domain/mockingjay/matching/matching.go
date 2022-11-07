package matching

import (
	"io"
	"net/http"
	"net/textproto"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"golang.org/x/exp/slices"
)

type Report struct {
	Matches         []RequestMatch     `json:"matches"`
	IncomingRequest mockingjay.Request `json:"incomingRequest"`
}

func NewReport(req *http.Request, endpoints mockingjay.Endpoints) Report {
	overallReport := Report{
		IncomingRequest: mockingjay.Request{
			Method:  req.Method,
			Path:    req.URL.String(),
			Headers: mockingjay.Headers(req.Header),
		},
	}
	reporter := newMatcher(req)
	for _, endpoint := range endpoints {
		overallReport.Matches = append(overallReport.Matches, reporter(endpoint))
	}

	return overallReport
}

func (m Report) FindMatchingResponse() (mockingjay.Response, bool) {
	i := slices.IndexFunc(m.Matches, func(r RequestMatch) bool {
		return r.Matched()
	})
	if i == -1 {
		return mockingjay.Response{}, false
	}
	return m.Matches[i].Response, true
}

func (m Report) HadMatch() bool {
	_, found := m.FindMatchingResponse()
	return found
}

type RequestMatch struct {
	Request  mockingjay.Request  `json:"request"`
	Response mockingjay.Response `json:"response"`
	Match    Match               `json:"match"`
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
			Request:  e.Request,
			Response: e.Response,
			Match: Match{
				Path:    e.Request.Path == req.URL.String(),
				Method:  e.Request.Method == req.Method,
				Headers: matchHeaders(e, req.Header),
				Body:    string(body) == e.Request.Body,
			},
		}
	}
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
