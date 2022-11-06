package endpoints

import (
	"net/http"
)

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

func (e Endpoints) GetMatchReport(req *http.Request) MatchReport {
	overallReport := MatchReport{
		IncomingRequest: Request{
			Method:  req.Method,
			Path:    req.URL.String(),
			Headers: Headers(req.Header),
		},
	}
	reporter := MatchReportFactory(req)
	for _, endpoint := range e.Endpoints {
		overallReport.Matches = append(overallReport.Matches, reporter(endpoint))
	}

	return overallReport
}
