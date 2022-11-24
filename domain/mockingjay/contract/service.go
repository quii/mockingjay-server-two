package contract

import (
	"net/http"

	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

type Service struct {
	httpClient *http.Client
}

func NewService(httpClient *http.Client) *Service {
	return &Service{httpClient: httpClient}
}

func (s Service) GetReports(endpoint http2.Endpoint) ([]Report, error) {
	var allReports []Report

	for _, cdc := range endpoint.CDCs {
		req := endpoint.Request.ToHTTPRequest(cdc.BaseURL)
		res, err := s.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		allReports = append(allReports, createReport(endpoint, res))
	}

	return allReports, nil
}

func createReport(endpoint http2.Endpoint, res *http.Response) Report {
	responseFromDownstream := http2.NewResponseFromHTTP(res)
	report := Report{
		Endpoint:               endpoint,
		ResponseFromDownstream: responseFromDownstream,
		Passed:                 IsResponseCompatible(responseFromDownstream, endpoint.Response),
	}
	return report
}

func IsResponseCompatible(got, want http2.Response) bool {
	return got.Status == want.Status && got.Body == want.Body //todo: match headers, re-use the other code probs
}
