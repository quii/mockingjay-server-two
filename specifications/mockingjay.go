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
	Configure(endpoints endpoints.Endpoints) error
}

type Client interface {
	Send(request endpoints.Request) (endpoints.Response, error)
}

func StubServerSpecification(t *testing.T, endpoints endpoints.Endpoints, mockingjay Mockingjay) {
	t.Run("mj can be configured with an endpoint, which can then be called by a client", func(t *testing.T) {
		assert.NoError(t, mockingjay.Configure(endpoints))

		for _, endpoint := range endpoints.Endpoints {
			t.Run(endpoint.Description, func(t *testing.T) {
				res, err := mockingjay.Send(endpoint.Request)
				assert.NoError(t, err)
				assertResponseMatches(t, endpoint.Response, res)
			})
		}
	})
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
