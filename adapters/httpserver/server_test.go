package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/drivers"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/specifications"
)

const (
	examplesDir = "../../specifications/examples/"
	fixturesDir = "../../specifications/testfixtures/"
)

func TestApp(t *testing.T) {
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir)
	assert.NoError(t, err)
	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir)
	assert.NoError(t, err)

	stubServer := httptest.NewServer(nil)
	adminServer := httptest.NewServer(nil)
	defer adminServer.Close()
	defer stubServer.Close()

	stubServerHandler, adminHandler := httpserver.New(
		matching.NewMockingjayStubServerService(examples),
		adminServer.URL,
	)

	stubServer.Config.Handler = stubServerHandler
	adminServer.Config.Handler = adminHandler

	client := &http.Client{Timeout: 2 * time.Second}
	driver := drivers.NewHTTPDriver(
		stubServer.URL,
		adminServer.URL,
		client,
	)
	webDriver := drivers.NewWebDriver(adminServer.URL, client)

	specifications.MockingjayStubServerSpec(t, driver, examples, fixtures)
	specifications.MockingjayAdmin(t, webDriver, examples)

	t.Run("view report", func(t *testing.T) {
		t.Run("404 if report doesn't exist", func(t *testing.T) {
			location := adminServer.URL + handlers.ReportsPath + "/" + uuid.New().String()
			_, err := driver.GetReport(location)

			assert.Error(t, err)
			assert.Equal(t, drivers.ErrReportNotFound{
				StatusCode: http.StatusNotFound,
				Location:   location,
			}, err.(drivers.ErrReportNotFound))
		})

		t.Run("404 if uuid wasnt valid", func(t *testing.T) {
			location := adminServer.URL + handlers.ReportsPath + "/" + "whatever"
			_, err := driver.GetReport(location)

			assert.Error(t, err)
			assert.Equal(t, drivers.ErrReportNotFound{
				StatusCode: http.StatusNotFound,
				Location:   location,
			}, err.(drivers.ErrReportNotFound))
		})
	})

	t.Run("put new configuration", func(t *testing.T) {
		t.Run("400 if you put a bad configuration", func(t *testing.T) {
			t.Run("invalid regex", func(t *testing.T) {
				assert.Error(t, driver.Configure(mockingjay.Endpoint{
					Description: "lala",
					Request: mockingjay.Request{
						Method:    http.MethodGet,
						RegexPath: "[", //invalid regex
						Path:      "/lol",
					},
				}))
			})

			t.Run("invalid structure", func(t *testing.T) {
				url := adminServer.URL + handlers.EndpointsPath
				req, _ := http.NewRequest(http.MethodPut, url, strings.NewReader("not valid JSON"))
				res, err := client.Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		})
	})
}
