package httpserver_test

import (
	"log"
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

	assert.NoError(t, examples.Compile())

	log.Println(examples[0].ID.String())

	service, err := matching.NewMockingjayStubServerService(nil)
	assert.NoError(t, err)

	stubServerHandler, adminHandler := httpserver.New(
		service,
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
	webDriver, cleanup := drivers.NewWebDriver(adminServer.URL, client, false)
	t.Cleanup(cleanup)

	specifications.MockingjayStubServerSpec(t, driver, examples, fixtures)

	t.Run("smaller ad-hoc example", func(t *testing.T) {
		endpoint := mockingjay.Endpoint{
			ID:          uuid.New(),
			Description: "Hello",
			Request: mockingjay.Request{
				Method:    http.MethodPost,
				RegexPath: "",
				Path:      "/todos/123",
				Headers: mockingjay.Headers{
					"accept": []string{"application/json"},
				},
				Body: "walk the dog",
			},
			Response: mockingjay.Response{
				Status: http.StatusNoContent,
				Body:   `{"task": "walk the dog"}`,
				Headers: mockingjay.Headers{
					"content-type": []string{"application/json"},
				},
			},
			CDCs: nil,
		}
		assert.NoError(t, webDriver.DeleteAllEndpoints())
		assert.NoError(t, webDriver.AddEndpoints(endpoint))
		endpoints, err := webDriver.GetEndpoints()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(endpoints))
		specifications.AssertEndpointEqual(t, endpoints[0], endpoint)
	})

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

		t.Run("404 if uuid wasn't valid", func(t *testing.T) {
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
				assert.Error(t, driver.AddEndpoints(mockingjay.Endpoint{
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
				req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader("not valid JSON"))
				req.Header.Add("content-type", "application/json")
				res, err := client.Do(req)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			})
		})
	})
}
