package specifications

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

type Configurer interface {
	Configure(endpoints endpoints.Endpoints) error
}

type Client interface {
	Do(request endpoints.Request) (endpoints.Response, error)
}

func GreetSpecification(t *testing.T, endpoints endpoints.Endpoints, configurer Configurer, client Client) {
	t.Run("mj can be configured with an endpoint, which can then be called by a client", func(t *testing.T) {
		assert.NoError(t, configurer.Configure(endpoints))

		for _, endpoint := range endpoints.Endpoints {
			t.Run(fmt.Sprintf("%q endpoint gets the correct response for its configured request", endpoint.Description), func(t *testing.T) {
				res, err := client.Do(endpoint.Request)
				assert.NoError(t, err)
				assert.Equal(t, endpoint.Response, res)
			})
		}
	})
}
