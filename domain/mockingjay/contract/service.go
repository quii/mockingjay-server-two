package contract

import (
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type Service struct {
	httpClient *http.Client
}

func NewService(httpClient *http.Client) *Service {
	return &Service{httpClient: httpClient}
}

func (s Service) GetReports(endpoint stub.Endpoint) ([]Report, error) {
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

func createReport(endpoint stub.Endpoint, res *http.Response) Report {
	got := stub.NewResponseFromHTTP(res)
	report := Report{
		Endpoint:               endpoint,
		ResponseFromDownstream: got,
		Passed:                 IsResponseCompatible(got, endpoint.Response),
	}
	return report
}
