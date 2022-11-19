package specifications

import (
	"fmt"
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"golang.org/x/exp/slices"
)

type Admin interface {
	GetReports() ([]matching.Report, error)
	DeleteReports() error

	AddEndpoints(endpoints ...mockingjay.Endpoint) error
	GetEndpoints() (mockingjay.Endpoints, error)
	DeleteEndpoint(uuid uuid.UUID) error
	DeleteEndpoints() error
}

type Client interface {
	Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error)
	//CheckEndpoints() ([]contract.Report, error) - wip
}

func MockingjayStubServerSpec(t *testing.T, admin Admin, client Client, examples mockingjay.Endpoints, testFixtures []mockingjay.TestFixture) {
	assert.NoError(t, admin.DeleteEndpoints())
	assert.NoError(t, admin.DeleteReports())

	t.Run("mj can be configured with request/response pairs (examples), which can then be called by a client with a request to get matching response", func(t *testing.T) {
		for _, endpoint := range examples {
			t.Run(endpoint.Description, func(t *testing.T) {
				assert.NoError(t, admin.AddEndpoints(endpoint))
				endpoints, err := admin.GetEndpoints()
				assert.NoError(t, err)
				assert.Equal(t, 1, len(endpoints))

				t.Cleanup(func() {
					assert.NoError(t, admin.DeleteEndpoint(endpoints[0].ID))
				})
				res, report, err := client.Send(endpoint.Request)
				assert.True(t, report.HadMatch, report)
				assert.NoError(t, err)
				assertResponseMatches(t, endpoint.Response, res)
			})
		}

		t.Run("a report of all requests made is available", func(t *testing.T) {
			reports, err := admin.GetReports()
			assert.NoError(t, err)
			t.Log(reports)
			assert.Equal(t, len(examples), len(reports))
		})
	})

	t.Run("mj test fixtures", func(t *testing.T) {
		for _, f := range testFixtures {
			t.Run(f.Endpoint.Description, func(t *testing.T) {
				t.Run("can configure mj on the fly with an endpoint", func(t *testing.T) {
					assert.NoError(t, admin.DeleteEndpoints())
					assert.NoError(t, admin.AddEndpoints(f.Endpoint))
					currentEndpoints, err := admin.GetEndpoints()
					assert.NoError(t, err)
					AssertEndpointsEqual(t, mockingjay.Endpoints{f.Endpoint}, currentEndpoints)
				})

				for _, request := range f.MatchingRequests {
					t.Run(request.Description, func(t *testing.T) {
						res, report, err := client.Send(request.Request)
						assert.NoError(t, err)
						assert.True(t, report.HadMatch, pp.Format(report))
						assertResponseMatches(t, f.Endpoint.Response, res)
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

func assertResponseMatches(t *testing.T, want, got mockingjay.Response) {
	t.Helper()
	assert.Equal(t, fudgeTheWhiteSpace(want.Body), fudgeTheWhiteSpace(got.Body), "body not equal")
	assert.Equal(t, want.Status, want.Status, "status not equal")

	for key, v := range want.Headers {
		for _, value := range v {
			i := slices.Index(got.Headers[key], value)
			assert.NotEqual(t, -1, i, fmt.Sprintf("%q not found in %v", value, got.Headers[key]))
		}
	}
}
