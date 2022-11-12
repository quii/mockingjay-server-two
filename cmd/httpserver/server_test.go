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

const (
	examplesDir = "../../examples/"
	fixturesDir = "../../testfixtures/"
)

func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	examples, err := mockingjay.NewEndpointsFromCue(examplesDir, os.DirFS(examplesDir))
	assert.NoError(t, err)
	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir, os.DirFS(fixturesDir))
	assert.NoError(t, err)

	var (
		stubServerPort   = "8080"
		configServerPort = "8081"
		driver           = httpserver.NewDriver(
			fmt.Sprintf("http://localhost:%s", stubServerPort),
			fmt.Sprintf("http://localhost:%s", configServerPort),
			&http.Client{
				Timeout: 1 * time.Second,
			},
		)
	)

	adapters.StartDockerServer(t, stubServerPort, configServerPort)
	specifications.MockingjaySpec(t, driver, examples, fixtures)
}
