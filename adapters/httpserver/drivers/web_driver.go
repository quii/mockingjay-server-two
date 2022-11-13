package drivers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type WebDriver struct {
	adminReportsURL   string
	adminEndpointsURL string
	client            *http.Client
	browser           *rod.Browser
}

func NewWebDriver(adminServerURL string, client *http.Client) *WebDriver {
	browser := rod.New().MustConnect()
	return &WebDriver{
		client:            client,
		browser:           browser,
		adminReportsURL:   adminServerURL + handlers.ReportsPath,
		adminEndpointsURL: adminServerURL + handlers.EndpointsPath,
	}
}

var (
	ErrNotImplemented = errors.New("not implemented")
)

func (d WebDriver) GetCurrentConfiguration() (mockingjay.Endpoints, error) {
	page := d.browser.MustPage(d.adminEndpointsURL)
	elements, err := page.Elements("tbody tr")
	if err != nil {
		return nil, err
	}
	var endpoints mockingjay.Endpoints
	for _, el := range elements {
		endpoints = append(endpoints, mockingjay.Endpoint{Description: el.MustText()})
	}
	return endpoints, nil
}

func (d WebDriver) Configure(es ...mockingjay.Endpoint) error {
	endpointJSON, err := json.Marshal(es)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, d.adminEndpointsURL, bytes.NewReader(endpointJSON))
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

	return nil
}

func (d WebDriver) Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error) {
	return mockingjay.Response{}, matching.Report{}, ErrNotImplemented
}

func (d WebDriver) GetReports() ([]matching.Report, error) {
	return nil, ErrNotImplemented
}

func (d WebDriver) GetReport(location string) (matching.Report, error) {
	return matching.Report{}, ErrNotImplemented
}
