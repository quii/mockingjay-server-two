package matching

import (
	"time"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Reports []Report

type Report struct {
	ID              uuid.UUID      `json:"ID"`
	HadMatch        bool           `json:"hadMatch"`
	IncomingRequest stub.Request   `json:"incomingRequest"`
	FailedMatches   []RequestMatch `json:"failed_matches"`
	SuccessfulMatch stub.Response  `json:"successfulMatch"`
	CreatedAt       time.Time      `json:"createdAt"`
}

func NewReport(req stub.Request, endpoints stub.Endpoints) Report {
	overallReport := Report{
		ID:              uuid.New(),
		IncomingRequest: req,
		CreatedAt:       time.Now().UTC(),
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
