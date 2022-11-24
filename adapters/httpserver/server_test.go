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
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
)

const (
	examplesDir = "../../specifications/examples/"
	fixturesDir = "../../specifications/stub_server_fixtures/"
)

func TestApp(t *testing.T) {
	stubServer, adminServer := startServers(t)

	service, err := mockingjay.NewService(nil, contract.NewService(&http.Client{}))
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

	examples, fixtures := mustLoadExamplesAndFixtures(t)

	t.Run("configuring with website", func(t *testing.T) {
		specifications.MockingjayStubServerSpec(t, webDriver, httpDriver, examples, fixtures)
	})

	t.Run("configuring with stub api", func(t *testing.T) {
		specifications.MockingjayStubServerSpec(t, httpDriver, httpDriver, examples, fixtures)
	})

	t.Run("consumer driven contracts", func(t *testing.T) {
		specifications.MockingjayConsumerDrivenContractSpec(t, httpDriver, httpDriver, examples)
	})

	t.Run("smaller ad-hoc example", func(t *testing.T) {
		endpoint := stub.Endpoint{
			ID:          uuid.New(),
			Description: "Hello",
			Request: stub.Request{
				Method:    http.MethodGet,
				RegexPath: "/happy-birthday/[a-z]+",
				Path:      "/happy-birthday/elodie",
				Headers: stub.Headers{
					"accept": []string{"application/json"},
				},
				Body: "walk the dog",
			},
			Response: stub.Response{
				Status: http.StatusOK,
				Body:   `{"msg": "happy birthday"}`,
				Headers: stub.Headers{
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
				assert.Error(t, httpDriver.AddEndpoints(stub.Endpoint{
					Description: "lala",
					Request: stub.Request{
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

func startServers(t *testing.T) (*httptest.Server, *httptest.Server) {
	stubServer := httptest.NewServer(nil)
	adminServer := httptest.NewServer(nil)
	t.Cleanup(stubServer.Close)
	t.Cleanup(adminServer.Close)
	return stubServer, adminServer
}

func mustLoadExamplesAndFixtures(t *testing.T) (stub.Endpoints, []mockingjay.TestFixture) {
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir)
	assert.NoError(t, err)
	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir)
	assert.NoError(t, err)
	return examples, fixtures
}
