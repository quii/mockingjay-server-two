package main_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/quii/mockingjay-server-two/adapters"
	"github.com/quii/mockingjay-server-two/adapters/config"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/drivers"
	"github.com/quii/mockingjay-server-two/specifications"
)

const specRoot = "../../specifications"

func TestMockingjayStubServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	adapters.StartMockingjayServer(t, config.DefaultStubServerPort, config.DefaultAdminServerPort)

	httpDriver := drivers.NewHTTPDriver(
		fmt.Sprintf("http://localhost:%s", config.DefaultStubServerPort),
		fmt.Sprintf("http://localhost:%s", config.DefaultAdminServerPort),
		&http.Client{
			Timeout: 1 * time.Second,
		},
	)
	specifications.MockingjayStubServerSpec(t, httpDriver, httpDriver, specRoot)
	specifications.MockingjayConsumerDrivenContractSpec(t, httpDriver, httpDriver, specRoot)
}
