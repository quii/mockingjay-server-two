package specifications

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/endpoints"
	"golang.org/x/exp/slices"
)

type Mockingjay interface {
	Configurer
	Client
}

type Configurer interface {
	Configure(endpoints ...endpoints.Endpoint) error
}

type Client interface {
	Send(request endpoints.Request) (endpoints.Response, endpoints.MatchReport, error)
}

type RequestDescription struct {
	Description string            `json:"description,omitempty"`
	Request     endpoints.Request `json:"request"`
}

type TestFixture struct {
	Endpoint            endpoints.Endpoint   `json:"endpoint"`
	MatchingRequests    []RequestDescription `json:"matchingRequests,omitempty"`
	NonMatchingRequests []RequestDescription `json:"nonMatchingRequests,omitempty"`
}

func MockingjaySpec(t *testing.T, mockingjay Mockingjay, examples endpoints.Endpoints, testFixtures []TestFixture) {
	t.Run("mj can be configured with request/response pairs, which can then be called by a client with a request to get matching response", func(t *testing.T) {
		assert.NoError(t, mockingjay.Configure(examples.Endpoints...))

		for _, endpoint := range examples.Endpoints {
			t.Run(endpoint.Description, func(t *testing.T) {
				res, report, err := mockingjay.Send(endpoint.Request)

				if !report.HadMatch() {
					t.Logf("Match report %#v", report)
				}

				assert.NoError(t, err)
				assertResponseMatches(t, endpoint.Response, res)
			})
		}
	})

	for _, f := range testFixtures {
		t.Run(f.Endpoint.Description, func(t *testing.T) {
			assert.NoError(t, mockingjay.Configure(f.Endpoint))

			for _, request := range f.MatchingRequests {
				t.Run(request.Description, func(t *testing.T) {
					res, _, err := mockingjay.Send(request.Request)
					assert.NoError(t, err)
					assertResponseMatches(t, f.Endpoint.Response, res)
				})
			}

			for _, request := range f.NonMatchingRequests {
				t.Run(request.Description, func(t *testing.T) {
					_, report, err := mockingjay.Send(request.Request)
					assert.NoError(t, err)
					assert.False(t, report.HadMatch())
				})
			}
		})
	}
}

func assertResponseMatches(t *testing.T, want, got endpoints.Response) {
	t.Helper()
	assert.Equal(t, want.Body, got.Body)
	assert.Equal(t, want.Status, want.Status)

	for key, v := range want.Headers {
		for _, value := range v {
			i := slices.Index(got.Headers[key], value)
			assert.NotEqual(t, -1, i, fmt.Sprintf("%q not found in %v", value, got.Headers[key]))
		}
	}
}
