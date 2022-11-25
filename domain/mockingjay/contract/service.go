package contract

import (
	"fmt"
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
		allReports = append(allReports, s.newReport(cdc, endpoint))
	}

	return allReports, nil
}

func (s Service) newReport(cdc stub.CDC, endpoint stub.Endpoint) Report {
	report := Report{
		Endpoint: endpoint,
		Ignore:   cdc.Ignore,
		URL:      cdc.BaseURL,
	}

	req := endpoint.Request.ToHTTPRequest(cdc.BaseURL)
	res, err := s.httpClient.Do(req)
	if err != nil {
		report.Errors = []string{fmt.Sprintf("could not reach %s, %s", req.URL, err)}
		return report
	}

	report.ResponseFromDownstream = stub.NewResponseFromHTTP(res)
	errors := IsResponseCompatible(report.ResponseFromDownstream, endpoint.Response)
	for _, err := range errors {
		report.Errors = append(report.Errors, err.Error())
	}
	return report
}
