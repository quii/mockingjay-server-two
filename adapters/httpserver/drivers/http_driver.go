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

type HTTPDriver struct {
	stubServerURL string
	reportsURL    string
	endpointsURL  string
	client        *http.Client
}

func (d HTTPDriver) DeleteReports() error {
	req, err := http.NewRequest(http.MethodDelete, d.reportsURL, nil)
	if err != nil {
		return err
	}
	res, err := d.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unepxected %d from %s", res.StatusCode, d.reportsURL)
	}
	return nil
}

func NewHTTPDriver(stubServerURL string, adminServerURL string, client *http.Client) *HTTPDriver {
	client.Transport = acceptJSONDecorator{transport: http.DefaultTransport}
	return &HTTPDriver{
		stubServerURL: stubServerURL,
		client:        client,
		reportsURL:    adminServerURL + handlers.ReportsPath,
		endpointsURL:  adminServerURL + handlers.EndpointsPath,
	}
}

func (d HTTPDriver) GetReports() ([]matching.Report, error) {
	var reports []matching.Report
	res, err := d.client.Get(d.reportsURL)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status of %d from %s", res.StatusCode, d.reportsURL)
	}

	err = json.NewDecoder(res.Body).Decode(&reports)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (d HTTPDriver) Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error) {
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

func (d HTTPDriver) GetReport(location string) (matching.Report, error) {
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
	if err := json.NewDecoder(res.Body).Decode(&matchReport); err != nil {
		return matching.Report{}, fmt.Errorf("could not decode response into reports %w", err)
	}
	return matchReport, nil
}

func (d HTTPDriver) AddEndpoints(es ...mockingjay.Endpoint) error {
	for _, e := range es {
		endpointJSON, err := json.Marshal(e)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, d.endpointsURL, bytes.NewReader(endpointJSON))
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
			return fmt.Errorf("got unexpected %d when trying to configure mj at %s", res.StatusCode, d.endpointsURL)
		}
	}

	return nil
}

func (d HTTPDriver) GetEndpoints() (mockingjay.Endpoints, error) {
	res, err := d.client.Get(d.endpointsURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", res.StatusCode, d.endpointsURL)
	}
	var endpoints mockingjay.Endpoints

	if err := json.NewDecoder(res.Body).Decode(&endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (d HTTPDriver) DeleteEndpoint(uuid uuid.UUID) error {
	url := d.endpointsURL + uuid.String()
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

func (d HTTPDriver) DeleteEndpoints() error {
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

type acceptJSONDecorator struct {
	transport http.RoundTripper
}

func (a acceptJSONDecorator) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	return a.transport.RoundTrip(req)
}
