package matching

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

type Reports []Report

type Report struct {
	ID              uuid.UUID      `json:"ID"`
	HadMatch        bool           `json:"hadMatch"`
	IncomingRequest http2.Request  `json:"incomingRequest"`
	FailedMatches   []RequestMatch `json:"failed_matches"`
	SuccessfulMatch http2.Response `json:"successfulMatch"`
	CreatedAt       time.Time      `json:"createdAt"`
}

func NewReport(req *http.Request, endpoints http2.Endpoints) Report {
	overallReport := Report{
		ID: uuid.New(),
		IncomingRequest: http2.Request{
			Method:  req.Method,
			Path:    req.URL.String(),
			Headers: http2.Headers(req.Header),
		},
		CreatedAt: time.Now().UTC(),
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

func SortReport(a, b Report) bool {
	return a.CreatedAt.Before(b.CreatedAt)
}
