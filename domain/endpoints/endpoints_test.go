package endpoints_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/endpoints"
	"github.com/quii/mockingjay-server-two/specifications"
)

func TestEndpointMatcher(t *testing.T) {
	const testFixturesDir = "../../examples/"
	fixtures, err := endpoints.NewEndpointsFromCue(testFixturesDir, os.DirFS(testFixturesDir))
	assert.NoError(t, err)

	driver := LocalDriver{}
	specifications.StubServerSpecification(t, fixtures, &driver, &driver)

	t.Run("it matches headers when the required header is not the first one", func(t *testing.T) {
		testEndpoints := endpoints.Endpoints{
			Endpoints: []endpoints.Endpoint{
				{Request: endpoints.Request{
					Method: http.MethodGet,
					Path:   "/",
					Headers: map[string][]string{
						"content-type": {"application/json"},
					},
				}},
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header = http.Header{
			"accept-encoding": {"gzip"},
			"Content-Type":    {"whatever", "application/json"},
		}

		report := testEndpoints.GetMatchReport(req)
		_, exists := report.FindMatchingResponse()
		assert.True(t, exists)
	})
}

type LocalDriver struct {
	endpoints endpoints.Endpoints
}

func (l *LocalDriver) Do(request endpoints.Request) (endpoints.Response, error) {
	report := l.endpoints.GetMatchReport(request.ToHTTPRequest(""))
	response, found := report.FindMatchingResponse()
	if !found {
		return endpoints.Response{}, fmt.Errorf("no request found for, match report %+v", report)
	}

	return response, nil
}

func (l *LocalDriver) Configure(endpoints endpoints.Endpoints) error {
	l.endpoints = endpoints
	return nil
}
