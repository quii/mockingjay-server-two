package specifications

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mj"
)

type Configurer interface {
	Configure(endpoints mj.Endpoints) error
}

type Client interface {
	Do(request mj.Request) (mj.Response, error)
}

func GreetSpecification(t *testing.T, configurer Configurer, client Client) {
	t.Run("mj can be configured with an endpoint, which can then be called by a client", func(t *testing.T) {
		var (
			helloWorldEndpoint = mj.Endpoint{
				Description: "Hello world endpoint",
				Request: mj.Request{
					Method: http.MethodGet,
					Path:   "/hello/world",
				},
				Response: mj.Response{
					Status: http.StatusOK,
					Body:   "Hello, world!",
				},
			}
			helloPepperEndpoint = mj.Endpoint{
				Description: "Hello pepper endpoint",
				Request: mj.Request{
					Method: http.MethodGet,
					Path:   "/hello/pepper",
				},
				Response: mj.Response{
					Status: http.StatusOK,
					Body:   "Hello, Pepper!",
				},
			}
			endpoints = []mj.Endpoint{helloWorldEndpoint, helloPepperEndpoint}
		)

		assert.NoError(t, configurer.Configure(mj.Endpoints{Endpoints: endpoints}))

		for _, endpoint := range endpoints {
			t.Run(fmt.Sprintf("%q endpoint gets the correct response for its configured request", endpoint.Description), func(t *testing.T) {
				res, err := client.Do(endpoint.Request)
				assert.NoError(t, err)
				assert.Equal(t, endpoint.Response, res)
			})
		}
	})
}
