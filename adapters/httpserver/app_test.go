package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/endpoints"
	"github.com/quii/mockingjay-server-two/specifications"
)

func TestApp(t *testing.T) {
	const testFixturesDir = "../../examples/"
	fixtures, err := endpoints.NewEndpointsFromCue(testFixturesDir, os.DirFS(testFixturesDir))
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

	specifications.StubServerSpecification(t, fixtures, driver)
}
