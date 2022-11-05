package endpoints

import (
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

type RequestMatch struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Matches  Matches  `json:"matches"`
}

func (r RequestMatch) Matched() bool {
	return r.Matches.Path && r.Matches.Method && r.Matches.Headers
}

type Matches struct {
	Path    bool `json:"path"`
	Method  bool `json:"method"`
	Headers bool `json:"headers"`
}

func MatchReportFactory(req *http.Request) func(Endpoint) RequestMatch {
	return func(e Endpoint) RequestMatch {
		return RequestMatch{
			Request:  e.Request,
			Response: e.Response,
			Matches: Matches{
				Path:    e.Request.Path == req.URL.String(),
				Method:  e.Request.Method == req.Method,
				Headers: matchHeaders(e, req.Header),
			},
		}
	}
}

func matchHeaders(e Endpoint, incomingHeaders http.Header) bool {
	headersMatch := len(e.Request.Headers) == 0

	for key, values := range e.Request.Headers {
		headerValues := incomingHeaders[textproto.CanonicalMIMEHeaderKey(key)]
		if len(headerValues) == 0 {
			continue
		}
		for _, vInReq := range headerValues {
			for _, vInEndpoint := range values {
				if vInReq == vInEndpoint {
					headersMatch = true
				}
			}
		}
	}
	return headersMatch
}
