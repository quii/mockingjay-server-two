package acceptance_tests_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/adapters"
	"github.com/quii/mockingjay-server-two/adapters/config"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/drivers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/specifications"
)

const (
	examplesDir = "../../../specifications/examples/"
	fixturesDir = "../../../specifications/testfixtures/"
)

func TestMockingjay(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	fixtures, err := mockingjay.NewFixturesFromCue(fixturesDir)
	assert.NoError(t, err)
	examples, err := mockingjay.NewEndpointsFromCue(examplesDir)
	assert.NoError(t, err)
	adapters.StartDockerServer(t, config.DefaultStubServerPort, config.DefaultAdminServerPort)

	t.Run("stub server", func(t *testing.T) {
		driver := drivers.NewHTTPDriver(
			fmt.Sprintf("http://localhost:%s", config.DefaultStubServerPort),
			fmt.Sprintf("http://localhost:%s", config.DefaultAdminServerPort),
			&http.Client{
				Timeout: 1 * time.Second,
			},
		)
		specifications.MockingjayStubServerSpec(t, driver, examples, fixtures)
	})

	t.Run("web interface", func(t *testing.T) {
		driver := drivers.NewWebDriver(
			fmt.Sprintf("http://localhost:%s", config.DefaultAdminServerPort),
			&http.Client{
				Timeout: 1 * time.Second,
			},
		)
		specifications.MockingjayAdmin(t, driver, examples)
	})
}
