package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Driver struct {
	stubServerURL     string
	adminReportsURL   string
	adminEndpointsURL string
	client            *http.Client
}

func NewHTTPDriver(stubServerURL string, adminServerURL string, client *http.Client) *Driver {
	return &Driver{
		stubServerURL:     stubServerURL,
		client:            client,
		adminReportsURL:   adminServerURL + handlers.ReportsPath,
		adminEndpointsURL: adminServerURL + handlers.EndpointsPath,
	}
}

func (d Driver) GetReports() ([]matching.Report, error) {
	var reports []matching.Report
	res, err := d.client.Get(d.adminReportsURL)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status of %d from %s", res.StatusCode, d.adminReportsURL)
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

	if res.Header.Get(handlers.HeaderMockingjayMatched) == "false" {
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
	req, err := http.NewRequest(http.MethodGet, location, nil)
	if err != nil {
		return matching.Report{}, err
	}
	req.Header.Set("Accept", "application/json")

	res, err := d.client.Do(req)
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
	if err := json.NewDecoder(res.Body).Decode(&matchReport); err != nil {
		return matching.Report{}, fmt.Errorf("could not decode response into reports %w", err)
	}
	return matchReport, nil
}

func (d Driver) AddEndpoints(es ...mockingjay.Endpoint) error {
	for _, e := range es {
		endpointJSON, err := json.Marshal(e)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, d.adminEndpointsURL, bytes.NewReader(endpointJSON))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			return err
		}

		res, err := d.client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusAccepted {
			return fmt.Errorf("got unexpected %d when trying to configure mj at %s", res.StatusCode, d.adminEndpointsURL)
		}
	}

	return nil
}

func (d Driver) GetEndpoints() (mockingjay.Endpoints, error) {
	req, _ := http.NewRequest(http.MethodGet, d.adminEndpointsURL, nil)
	req.Header.Set("Accept", "application/json")

	res, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", res.StatusCode, d.adminEndpointsURL)
	}
	var endpoints mockingjay.Endpoints

	if err := json.NewDecoder(res.Body).Decode(&endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (d Driver) DeleteEndpoint(uuid uuid.UUID) error {
	url := d.adminEndpointsURL + uuid.String()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	res, err := d.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unpexected %d from %s", res.StatusCode, url)
	}
	return nil
}

func (d Driver) DeleteAllEndpoints() error {
	endpoints, err := d.GetEndpoints()
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		if err := d.DeleteEndpoint(endpoint.ID); err != nil {
			return err
		}
	}
	return nil
}

type ErrReportNotFound struct {
	StatusCode int
	Location   string
}

func (e ErrReportNotFound) Error() string {
	return fmt.Sprintf("unexpected %d from %s", e.StatusCode, e.Location)
}
