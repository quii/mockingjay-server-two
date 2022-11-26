package stub_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func TestEndpoint_Compile(t *testing.T) {
	t.Run("sets the loaded time", func(t *testing.T) {
		beforeCompile := time.Now()
		endpoint := stub.Endpoint{
			LoadedAt: time.Time{},
		}
		assert.NoError(t, endpoint.Compile())
		assert.True(t, endpoint.LoadedAt.After(beforeCompile))
	})

	t.Run("sets headers into canonical format", func(t *testing.T) {
		endpoint := stub.Endpoint{Request: stub.Request{Headers: map[string][]string{
			"aCCEPt": {"application/json"},
		}}, Response: stub.Response{Headers: map[string][]string{
			"cOnTenT-tyPE": {"application/json"},
		}}}

		assert.NoError(t, endpoint.Compile())
		assert.Equal(t, endpoint.Request.Headers["Accept"], []string{"application/json"})
		assert.Equal(t, endpoint.Response.Headers["Content-Type"], []string{"application/json"})
	})

	t.Run("compile all", func(t *testing.T) {
		beforeCompile := time.Now()

		endpoints := stub.Endpoints{
			stub.Endpoint{},
			stub.Endpoint{},
		}

		assert.NoError(t, endpoints.Compile())
		for _, endpoint := range endpoints {
			assert.True(t, endpoint.LoadedAt.After(beforeCompile))
		}
	})
}
