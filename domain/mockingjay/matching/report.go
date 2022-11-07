package matching

import (
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
)

type Report struct {
	FailedMatches   []RequestMatch      `json:"failed_matches"`
	SuccessfulMatch mockingjay.Response `json:"successfulMatch"`
	IncomingRequest mockingjay.Request  `json:"incomingRequest"`
	HadMatch        bool                `json:"hadMatch"`
}

func NewReport(req *http.Request, endpoints mockingjay.Endpoints) Report {
	overallReport := Report{
		IncomingRequest: mockingjay.Request{
			Method:  req.Method,
			Path:    req.URL.String(),
			Headers: mockingjay.Headers(req.Header),
		},
	}
	matcher := newMatcher(req)
	for _, endpoint := range endpoints {
		match := matcher(endpoint)
		if match.Matched() {
			overallReport.SuccessfulMatch = match.Endpoint.Response
			overallReport.HadMatch = true
		} else {
			overallReport.FailedMatches = append(overallReport.FailedMatches, match)
		}
	}

	return overallReport
}
