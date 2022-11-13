package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Driver struct {
	stubServerURL      string
	client             *http.Client
	matchReportURL     string
	configEndpointsURL string
	adminServerURL     string
}

func NewDriver(stubServerURL string, adminServerURL string, client *http.Client) *Driver {
	return &Driver{
		stubServerURL:      stubServerURL,
		adminServerURL:     adminServerURL,
		client:             client,
		matchReportURL:     adminServerURL + ReportsPath,
		configEndpointsURL: adminServerURL + EndpointsPath,
	}
}

func (d Driver) GetReports() ([]matching.Report, error) {
	var reports []matching.Report
	res, err := d.client.Get(d.matchReportURL)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status of %d from %s", res.StatusCode, d.matchReportURL)
	}

	err = json.NewDecoder(res.Body).Decode(&reports)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (d Driver) Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error) {
	req := request.ToHTTPRequest(d.stubServerURL)

	res, err := d.client.Do(req)
	if err != nil {
		return mockingjay.Response{}, matching.Report{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mockingjay.Response{}, matching.Report{}, err
	}

	if res.Header.Get(HeaderMockingjayMatched) == "false" {
		report, err := d.GetReport(res.Header.Get("location"))
		if err != nil {
			return mockingjay.Response{}, report, err
		}
		return mockingjay.Response{}, report, nil
	}

	return mockingjay.Response{
		Status:  res.StatusCode,
		Body:    string(body),
		Headers: mockingjay.Headers(res.Header),
	}, matching.Report{HadMatch: true}, nil
}

func (d Driver) GetReport(location string) (matching.Report, error) {
	var matchReport matching.Report
	res, err := d.client.Get(location)
	if err != nil {
		return matchReport, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return matchReport, ErrReportNotFound{
			StatusCode: res.StatusCode,
			Location:   location,
		}
	}
	_ = json.NewDecoder(res.Body).Decode(&matchReport)
	return matching.Report{}, nil
}

func (d Driver) Configure(es ...mockingjay.Endpoint) error {
	endpointJSON, err := json.Marshal(es)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, d.configEndpointsURL, bytes.NewReader(endpointJSON))
	if err != nil {
		return err
	}

	res, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("got unexpected %d when trying to configure mj at %s", res.StatusCode, d.configEndpointsURL)
	}

	return nil
}

func (d Driver) GetCurrentConfiguration() (mockingjay.Endpoints, error) {
	res, err := d.client.Get(d.configEndpointsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", res.StatusCode, d.configEndpointsURL)
	}
	var endpoints mockingjay.Endpoints

	if err := json.NewDecoder(res.Body).Decode(&endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

type ErrReportNotFound struct {
	StatusCode int
	Location   string
}

func (e ErrReportNotFound) Error() string {
	return fmt.Sprintf("unexpected %d from %s", e.StatusCode, e.Location)
}
