package mockingjay_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/adapters/httpserver/handlers"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

func TestService_CheckEndpoints(t *testing.T) {
	httpClient := &http.Client{}

	service, _ := mockingjay.NewService(nil, contract.NewService(httpClient))
	downstreamService := httptest.NewServer(handlers.NewStubHandler(service, "n/a"))
	t.Cleanup(downstreamService.Close)

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

	assert.NoError(t, service.Endpoints().Create(endpoint.ID, endpoint))

	t.Run("it can check the cdcs", func(t *testing.T) {
		reports, err := service.CheckEndpoints()
		assert.NoError(t, err)
		assert.Equal(t, len(endpoint.CDCs), len(reports))
		assert.True(t, reports[0].Passed, pp.Format(reports[0]))
	})

}
