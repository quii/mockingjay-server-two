package drivers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		getRequestField := func(field string) string {
			return el.MustElement(fmt.Sprintf(`td[data-request-field=%s]`, field)).MustText()
		}
		getResponseField := func(field string) string {
			return el.MustElement(fmt.Sprintf(`td[data-response-field=%s]`, field)).MustText()
		}
		statusText := getResponseField("status")
		statusCode, err := strconv.Atoi(statusText)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, mockingjay.Endpoint{
			Description: el.MustElement(`td[data-field=description`).MustText(),
			Request: mockingjay.Request{
				Method:    getRequestField("method"),
				RegexPath: getRequestField("regexPath"),
				Path:      getRequestField("path"),
				Headers:   extractHeadersFromMarkup(el, `td[data-request-field=headers] dl`),
				Body:      getRequestField("body"),
			},
			Response: mockingjay.Response{
				Status:  statusCode,
				Body:    getResponseField("body"),
				Headers: extractHeadersFromMarkup(el, `td[data-response-field=headers] dl`),
			},
		})
	}
	return endpoints, nil
}

func extractHeadersFromMarkup(el *rod.Element, selector string) mockingjay.Headers {
	var requestHeaders mockingjay.Headers
	dl := el.MustElement(selector)
	listItems := dl.MustElements("*")

	if len(listItems) > 0 {
		requestHeaders = make(mockingjay.Headers)
		currentKey := ""
		for _, item := range listItems {
			if item.String() == "<dt>" {
				currentKey = item.MustText()
				continue
			}
			requestHeaders[currentKey] = append(requestHeaders[currentKey], item.MustText())
		}
	}
	return requestHeaders
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
