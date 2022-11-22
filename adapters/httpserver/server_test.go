package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters/config"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/drivers"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/specifications"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
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

	assert.NoError(t, examples.Compile())
	service, err := mockingjay.NewStubService(nil)
	assert.NoError(t, err)

	stubServerHandler, adminHandler, err := httpserver.New(
		service,
		adminServer.URL,
		config.DevModeOn,
	)
	assert.NoError(t, err)

	stubServer.Config.Handler = stubServerHandler
	adminServer.Config.Handler = adminHandler

	client := &http.Client{Timeout: 2 * time.Second}
	httpDriver := drivers.NewHTTPDriver(
		stubServer.URL,
		adminServer.URL,
		client,
	)
	webDriver, cleanup := drivers.NewWebDriver(adminServer.URL, client, false)
	t.Cleanup(cleanup)

	t.Run("configuring with website", func(t *testing.T) {
		specifications.MockingjayStubServerSpec(t, webDriver, httpDriver, examples, fixtures)
	})

	t.Run("configuring with http api", func(t *testing.T) {
		specifications.MockingjayStubServerSpec(t, httpDriver, httpDriver, examples, fixtures)
	})

	t.Run("smaller ad-hoc example", func(t *testing.T) {
		endpoint := http2.Endpoint{
			ID:          uuid.New(),
			Description: "Hello",
			Request: http2.Request{
				Method:    http.MethodGet,
				RegexPath: "/happy-birthday/[a-z]+",
				Path:      "/happy-birthday/elodie",
				Headers: http2.Headers{
					"accept": []string{"application/json"},
				},
				Body: "walk the dog",
			},
			Response: http2.Response{
				Status: http.StatusOK,
				Body:   `{"msg": "happy birthday"}`,
				Headers: http2.Headers{
					"content-type": []string{"application/json"},
				},
			},
			CDCs: nil,
		}
		assert.NoError(t, endpoint.Compile())
		usecases.StubServer{
			Admin:  webDriver,
			Client: httpDriver,
		}.Test(t, endpoint)
	})

	t.Run("view report", func(t *testing.T) {
		t.Run("404 if report doesn't exist", func(t *testing.T) {
			location := adminServer.URL + handlers.ReportsPath + "/" + uuid.New().String()
			_, err := httpDriver.GetReport(location)

			assert.Error(t, err)
			assert.Equal(t, drivers.ErrReportNotFound{
				StatusCode: http.StatusNotFound,
				Location:   location,
			}, err.(drivers.ErrReportNotFound))
		})

		t.Run("404 if uuid wasn't valid", func(t *testing.T) {
			location := adminServer.URL + handlers.ReportsPath + "/" + "whatever"
			_, err := httpDriver.GetReport(location)

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
				assert.Error(t, httpDriver.AddEndpoints(http2.Endpoint{
					Description: "lala",
					Request: http2.Request{
						Method:    http.MethodGet,
						RegexPath: "[", //invalid regex
						Path:      "/lol",
					},
				}))
			})

			t.Run("invalid structure", func(t *testing.T) {
				url := adminServer.URL + handlers.EndpointsPath
				req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader("not valid JSON"))
				req.Header.Add("content-type", "application/json")
				res, err := client.Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		})
	})
}
