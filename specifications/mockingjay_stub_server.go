package specifications

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Admin interface {
	GetReports() ([]matching.Report, error)
	DeleteReports() error

	AddEndpoints(endpoints ...http.Endpoint) error
	GetEndpoints() (http.Endpoints, error)
	DeleteEndpoint(uuid uuid.UUID) error
	DeleteEndpoints() error
}

type Client interface {
	Send(request http.Request) (http.Response, matching.Report, error)
	//CheckEndpoints() ([]contract.Report, error) - wip
}

func MockingjayStubServerSpec(t *testing.T, admin Admin, client Client, examples http.Endpoints, testFixtures []mockingjay.TestFixture) {
	assert.NoError(t, admin.DeleteEndpoints())
	assert.NoError(t, admin.DeleteReports())

	stubServerUseCase := StubServerUseCase{
		admin:  admin,
		client: client,
	}

	t.Run("mj can be configured with request/response pairs (examples), which can then be called by a client with a request to get matching response", func(t *testing.T) {
		for _, example := range examples {
			stubServerUseCase.Test(t, example)
		}
	})

	t.Run("mj test fixtures", func(t *testing.T) {
		for _, f := range testFixtures {
			t.Run(f.Endpoint.Description, func(t *testing.T) {
				t.Run("can configure mj on the fly with an endpoint", func(t *testing.T) {
					assert.NoError(t, admin.DeleteEndpoints())
					assert.NoError(t, admin.AddEndpoints(f.Endpoint))
					currentEndpoints, err := admin.GetEndpoints()
					assert.NoError(t, err)
					AssertEndpointsEqual(t, http.Endpoints{f.Endpoint}, currentEndpoints)
				})

				for _, request := range f.MatchingRequests {
					t.Run(request.Description, func(t *testing.T) {
						res, report, err := client.Send(request.Request)
						assert.NoError(t, err)
						assert.True(t, report.HadMatch, pp.Format(report))
						AssertResponseMatches(t, f.Endpoint.Response, res)
					})
				}

				for _, request := range f.NonMatchingRequests {
					t.Run(request.Description, func(t *testing.T) {
						_, report, err := client.Send(request.Request)
						assert.NoError(t, err)
						assert.False(t, report.HadMatch)
					})
				}
			})
		}
	})
}
