package endpoints

import (
	"io"
	"net/http"
	"net/textproto"

	"golang.org/x/exp/slices"
)

type MatchReport struct {
	Matches         []RequestMatch `json:"matches"`
	IncomingRequest Request        `json:"incomingRequest"`
}

func (m MatchReport) FindMatchingResponse() (Response, bool) {
	i := slices.IndexFunc(m.Matches, func(r RequestMatch) bool {
		return r.Matched()
	})
	if i == -1 {
		return Response{}, false
	}
	return m.Matches[i].Response, true
}

func (m MatchReport) HadMatch() bool {
	_, found := m.FindMatchingResponse()
	return found
}

type RequestMatch struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Matches  Matches  `json:"matches"`
}

func (r RequestMatch) Matched() bool {
	return r.Matches.Path && r.Matches.Method && r.Matches.Headers && r.Matches.Body
}

type Matches struct {
	Path    bool `json:"path"`
	Method  bool `json:"method"`
	Headers bool `json:"headers"`
	Body    bool `json:"body"`
}

func MatchReportFactory(req *http.Request) func(Endpoint) RequestMatch {
	var body []byte
	if req.Body != nil {
		defer req.Body.Close()
		body, _ = io.ReadAll(req.Body)
	}

	return func(e Endpoint) RequestMatch {
		return RequestMatch{
			Request:  e.Request,
			Response: e.Response,
			Matches: Matches{
				Path:    e.Request.Path == req.URL.String(),
				Method:  e.Request.Method == req.Method,
				Headers: matchHeaders(e, req.Header),
				Body:    string(body) == e.Request.Body,
			},
		}
	}
}

func matchHeaders(e Endpoint, incomingHeaders http.Header) bool {
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
