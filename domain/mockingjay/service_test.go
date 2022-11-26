package mockingjay_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/specifications"
)

const specRoot = "../../specifications"

func TestService_CheckEndpoints(t *testing.T) {
	service := mockingjay.NewService(contract.NewService(&http.Client{}))
	downstreamService := httptest.NewServer(handlers.NewStubHandler(service, "n/a"))
	t.Cleanup(downstreamService.Close)
	driver := mockingjay.NewDriver(service)

	specifications.MockingjayStubServerSpec(t, driver, driver, specRoot)
	specifications.MockingjayConsumerDrivenContractSpec(t, driver, driver, specRoot)
}
