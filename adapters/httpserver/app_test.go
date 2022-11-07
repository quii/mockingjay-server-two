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
	const examplesDir = "../../examples/"
	examples, err := endpoints.NewEndpointsFromCue(examplesDir, os.DirFS(examplesDir))
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

	fixture := specifications.TestFixture{
		Endpoint: endpoints.Endpoint{
			Description: "This will be loaded from file later",
			Request: endpoints.Request{
				Method: http.MethodGet,
				Path:   "/",
				Headers: endpoints.Headers{
					"Accept": {"application/xml"},
				},
			},
			Response: endpoints.Response{
				Status: http.StatusOK,
				Body:   `<hello>World</hello>`,
				Headers: endpoints.Headers{
					"Content-Type": {"application/xml"},
				},
			},
		},
		MatchingRequests: []specifications.RequestDescription{
			{
				Description: "Works even though the header is not the first one",
				Request: endpoints.Request{
					Method: http.MethodGet,
					Path:   "/",
					Headers: endpoints.Headers{
						"Accept": {"text/html", "application/xml"},
					},
				},
			},
		},
		NonMatchingRequests: []specifications.RequestDescription{
			{
				Description: "Doesn't match as it has the wrong header",
				Request: endpoints.Request{
					Method: http.MethodGet,
					Path:   "/",
					Headers: endpoints.Headers{
						"Accept": {"text/html"},
					},
				},
			},
		},
	}

	specifications.MockingjaySpec(t, driver, examples, []specifications.TestFixture{fixture})
}
