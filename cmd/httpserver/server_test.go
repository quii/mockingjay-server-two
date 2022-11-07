package main_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/adapters"
	"github.com/quii/mockingjay-server-two/adapters/httpserver"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/specifications"
)

func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	const testFixturesDir = "../../examples/"
	examples, err := mockingjay.NewEndpointsFromCue(testFixturesDir, os.DirFS(testFixturesDir))
	assert.NoError(t, err)

	var (
		stubServerPort   = "8080"
		configServerPort = "8081"
		driver           = httpserver.Driver{
			StubServerURL:   fmt.Sprintf("http://localhost:%s", stubServerPort),
			ConfigServerURL: fmt.Sprintf("http://localhost:%s", configServerPort),
			Client: &http.Client{
				Timeout: 1 * time.Second,
			},
		}
	)

	adapters.StartDockerServer(t, stubServerPort, configServerPort)
	specifications.MockingjaySpec(t, driver, examples, nil)
}
