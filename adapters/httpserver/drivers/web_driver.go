package drivers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/google/uuid"
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
			SlowMotion(100 * time.Millisecond).
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

func (d WebDriver) GetEndpoints() (mockingjay.Endpoints, error) {
	var endpoints mockingjay.Endpoints

	endpointsPage, err := d.browser.Page(proto.TargetCreateTarget{URL: d.adminEndpointsURL})
	endpointsPage.MustWaitNavigation()
	endpointsPage.MustElement("table") // force the thing to check for a table

	if err != nil {
		return nil, fmt.Errorf("unable to reach %s, %w", d.adminEndpointsURL, err)
	}
	elements, err := endpointsPage.Elements(".endpoint")
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

func (d WebDriver) AddEndpoints(es ...mockingjay.Endpoint) error {
	page := d.browser.MustPage(d.adminEndpointsURL)
	for _, endpoint := range es {
		form, err := page.Element("form")
		if err != nil {
			return fmt.Errorf("couldn't find form in page to enter endpoint %w", err)
		}

		if err := pageobjects.InsertEndpoint(form, endpoint); err != nil {
			return err
		}
		page.MustWaitNavigation()
	}

	return nil
}

func (d WebDriver) DeleteEndpoint(uuid uuid.UUID) error {
	page := d.browser.MustPage(d.adminEndpointsURL)
	rowToDelete, err := page.Element(fmt.Sprintf(`*[data-id="%s"]`, uuid.String()))
	if err != nil {
		return err
	}
	rowToDelete.MustElement("button").MustClick()
	return nil
}

func (d WebDriver) DeleteAllEndpoints() error {
	endpoints, err := d.GetEndpoints()
	if err != nil {
		return err
	}
	log.Println("going to delete", len(endpoints))
	for _, endpoint := range endpoints {
		log.Println("deleting", endpoint.ID.String())
		if err := d.DeleteEndpoint(endpoint.ID); err != nil {
			return err
		}
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
