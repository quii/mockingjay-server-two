package main_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/adapters"
	"github.com/quii/mockingjay-server-two/adapters/config"
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

	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir)
	assert.NoError(t, err)
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir)
	assert.NoError(t, err)

	t.Run("loading configuration via admin server", func(t *testing.T) {
		driver := httpserver.NewDriver(
			fmt.Sprintf("http://localhost:%s", config.DefaultStubServerPort),
			fmt.Sprintf("http://localhost:%s", config.DefaultAdminServerPort),
			&http.Client{
				Timeout: 1 * time.Second,
			},
		)

		adapters.StartDockerServer(t, config.DefaultStubServerPort, config.DefaultAdminServerPort)
		specifications.MockingjaySpec(t, driver, examples, fixtures)
	})
}
