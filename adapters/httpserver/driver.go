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

/*
notes
- we can have mj listen on two ports, one for mj mgmt and the other for the stub server so we can configure without conflict. for now though, just the one
*/

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
		configEndpointsURL: adminServerURL + ConfigEndpointsPath,
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
	var matchReport matching.Report

	req := request.ToHTTPRequest(d.stubServerURL)

	res, err := d.client.Do(req)
	if err != nil {
		return mockingjay.Response{}, matchReport, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return mockingjay.Response{}, matchReport, err
	}

	if res.Header.Get(HeaderMockingjayMatched) == "false" {
		reportURL := d.adminServerURL + res.Header.Get("location")
		res, err := d.client.Get(reportURL)
		if err != nil {
			return mockingjay.Response{}, matchReport, err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return mockingjay.Response{}, matchReport, fmt.Errorf("unexpected %d from %s", res.StatusCode, reportURL)
		}
		json.NewDecoder(res.Body).Decode(&matchReport)
		return mockingjay.Response{}, matchReport, nil
	}

	matchReport.HadMatch = true

	return mockingjay.Response{
		Status:  res.StatusCode,
		Body:    string(body),
		Headers: mockingjay.Headers(res.Header),
	}, matchReport, nil
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
