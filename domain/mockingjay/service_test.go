package mockingjay_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
)

func TestService_CheckEndpoints(t *testing.T) {
	httpClient := &http.Client{}

	service := mockingjay.NewService(contract.NewService(httpClient))
	downstreamService := httptest.NewServer(handlers.NewStubHandler(service, "n/a"))
	t.Cleanup(downstreamService.Close)
	driver := mockingjay.NewDriver(service)

	endpoint := stub.Endpoint{
		ID: uuid.New(),
		Request: stub.Request{
			Method: http.MethodGet,
			Path:   "/",
		},
		Response: stub.Response{
			Status: http.StatusOK,
			Body:   "Hello, world",
		},
		CDCs: []stub.CDC{
			{
				BaseURL:   downstreamService.URL,
				Retries:   0,
				TimeoutMS: 0,
			},
		},
	}

	usecases.ConsumerDrivenContract{
		Admin:  driver,
		Client: driver,
	}.Test(t, endpoint)

	usecases.StubServer{
		Admin:  driver,
		Client: driver,
	}.Test(t, endpoint)
}
