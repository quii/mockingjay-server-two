package specifications

import (
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/config"
)

type Configurer interface {
	Configure(endpoints config.Endpoints) error
}

type Client interface {
	Do(request config.Request) (config.Response, error)
}

func GreetSpecification(t *testing.T, configurer Configurer, client Client) {
	t.Run("mj can be configured with an endpoint, which can then be called by a client", func(t *testing.T) {
		var (
			endpoint = config.Endpoint{
				Description: "Hello world endpoint",
				Request: config.Request{
					Method: http.MethodGet,
					Path:   "/hello/world",
				},
				Response: config.Response{
					Status: http.StatusOK,
					Body:   "Hello, world!",
				},
			}
			endpoints = []config.Endpoint{endpoint}
		)

		assert.NoError(t, configurer.Configure(config.Endpoints{Endpoints: endpoints}))

		res, err := client.Do(endpoint.Request)
		assert.NoError(t, err)
		assert.Equal(t, endpoint.Response, res)
	})
}
