package drivers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/drivers/internal/pageobjects"
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

var (
	ErrNotImplemented = errors.New("not implemented")
)

func NewWebDriver(adminServerURL string, client *http.Client, debug bool) (*WebDriver, func()) {
	var browser *rod.Browser
	var cleanup func()

	if debug {
		l := launcher.New().
			Headless(false).
			Devtools(true)

		cleanup = l.Cleanup

		url := l.MustLaunch()

		browser = rod.New().
			ControlURL(url).
			Trace(true).
			SlowMotion(500 * time.Millisecond).
			MustConnect()
	} else {
		browser = rod.New().MustConnect()
		cleanup = func() {}
	}

	return &WebDriver{
		client:            client,
		browser:           browser,
		adminReportsURL:   adminServerURL + handlers.ReportsPath,
		adminEndpointsURL: adminServerURL + handlers.EndpointsPath,
	}, cleanup
}

func (d WebDriver) GetCurrentConfiguration() (mockingjay.Endpoints, error) {
	var endpoints mockingjay.Endpoints

	elements, err := d.browser.MustPage(d.adminEndpointsURL).Elements(".endpoint")
	if err != nil {
		return nil, err
	}

	for _, el := range elements {
		endpoint, err := pageobjects.EndpointFromMarkup(el)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

func (d WebDriver) Configure(es ...mockingjay.Endpoint) error {
	page := d.browser.MustPage(d.adminEndpointsURL)
	for _, endpoint := range es {
		page.MustElement(`*[name="description"]`).MustInput(endpoint.Description)

		page.MustElement(`*[name="path"]`).MustInput(endpoint.Request.Path)
		page.MustElement(`*[name="regexpath"]`).MustInput(endpoint.Request.RegexPath)
		page.MustElement(`*[name="method"]`).MustSelect(endpoint.Request.Method)
		page.MustElement(`*[name="request.body"]`).MustInput(endpoint.Request.Body)

		for k, v := range endpoint.Request.Headers {
			page.MustElement(`*[name="request.header.name"]`).MustInput(k)
			page.MustElement(`*[name="request.header.values"]`).MustInput(strings.Join(v, "; "))
		}

		page.MustElement(`*[name="status"]`).MustInput(fmt.Sprintf("%d", endpoint.Response.Status))
		page.MustElement(`*[name="response.body"]`).MustInput(endpoint.Response.Body)

		for k, v := range endpoint.Request.Headers {
			page.MustElement(`*[name="response.header.name"]`).MustInput(k)
			page.MustElement(`*[name="response.header.values"]`).MustInput(strings.Join(v, "; "))
		}

		submitButton, err := page.Element(`#submit`)
		if err != nil {
			return err
		}
		if err := submitButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
			return err
		}
		page.MustWaitNavigation()
	}

	return nil
}

func (d WebDriver) Reset() error {
	req, err := http.NewRequest(http.MethodDelete, d.adminEndpointsURL, nil)
	if err != nil {
		return err
	}
	res, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("got unexpected %d when trying to reset mj at %s", res.StatusCode, d.adminEndpointsURL)
	}

	return nil
}

func (d WebDriver) Send(_ mockingjay.Request) (mockingjay.Response, matching.Report, error) {
	return mockingjay.Response{}, matching.Report{}, ErrNotImplemented
}

func (d WebDriver) GetReports() ([]matching.Report, error) {
	return nil, ErrNotImplemented
}

func (d WebDriver) GetReport(_ string) (matching.Report, error) {
	return matching.Report{}, ErrNotImplemented
}
