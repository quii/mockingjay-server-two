package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/specifications"
)

const (
	examplesDir = "../../examples/"
	fixturesDir = "../../testfixtures/"
)

func TestApp(t *testing.T) {
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir)
	assert.NoError(t, err)
	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir)
	assert.NoError(t, err)

	app := httpserver.New(examples, "")
	stubServer := httptest.NewServer(http.HandlerFunc(app.StubHandler))
	adminServer := httptest.NewServer(app.AdminRouter)
	defer adminServer.Close()
	defer stubServer.Close()
	app.AdminBaseURL = adminServer.URL

	driver := httpserver.NewDriver(
		stubServer.URL,
		adminServer.URL,
		&http.Client{},
	)

	specifications.MockingjaySpec(t, driver, examples, fixtures)

	t.Run("view report", func(t *testing.T) {
		t.Run("404 if report doesn't exist", func(t *testing.T) {
			_, err := driver.GetReport(httpserver.ReportsPath + "/whatever")
			assert.Error(t, err)
		})
	})
}
