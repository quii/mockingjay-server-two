package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir, os.DirFS(examplesDir))
	assert.NoError(t, err)
	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir, os.DirFS(fixturesDir))
	assert.NoError(t, err)

	app := new(httpserver.App)
	stubServer := httptest.NewServer(http.HandlerFunc(app.StubHandler))
	configServer := httptest.NewServer(http.HandlerFunc(app.ConfigHandler))
	defer configServer.Close()
	defer stubServer.Close()

	driver := httpserver.Driver{
		StubServerURL:   stubServer.URL,
		ConfigServerURL: configServer.URL,
		Client:          &http.Client{},
	}

	specifications.MockingjaySpec(t, driver, examples, fixtures)
}
