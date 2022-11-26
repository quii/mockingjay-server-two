package stub_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func TestEndpoint_Compile(t *testing.T) {
	t.Run("when compiling, sets the loaded time", func(t *testing.T) {
		beforeCompile := time.Now()
		endpoint := stub.Endpoint{
			LoadedAt: time.Time{},
		}
		assert.NoError(t, endpoint.Compile())
		assert.True(t, endpoint.LoadedAt.After(beforeCompile))
	})
}
