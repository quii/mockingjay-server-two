package specifications

import (
	"fmt"
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"golang.org/x/exp/slices"
)

type Mockingjay interface {
	Admin
	Client
}

type Client interface {
	Send(request mockingjay.Request) (mockingjay.Response, matching.Report, error)
	//CheckEndpoints() ([]contract.Report, error) - wip
}

func MockingjayStubServerSpec(t *testing.T, mj Mockingjay, examples mockingjay.Endpoints, testFixtures []mockingjay.TestFixture) {
	t.Run("mj can be pre-configured with request/response pairs (examples), which can then be called by a client with a request to get matching response", func(t *testing.T) {
		for _, endpoint := range examples {
			t.Run(endpoint.Description, func(t *testing.T) {
				res, report, err := mj.Send(endpoint.Request)

				if !report.HadMatch {
					t.Logf("Match report %#v", report)
				}

				assert.NoError(t, err)
				assertResponseMatches(t, endpoint.Response, res)
			})
		}

		t.Run("a report of all requests made is available", func(t *testing.T) {
			reports, err := mj.GetReports()
			assert.NoError(t, err)
			assert.Equal(t, len(examples), len(reports))
		})
	})

	t.Run("mj test fixtures", func(t *testing.T) {
		for _, f := range testFixtures {
			t.Run(f.Endpoint.Description, func(t *testing.T) {
				assert.NoError(t, mj.Configure(f.Endpoint))
				currentEndpoints, err := mj.GetCurrentConfiguration()
				assert.NoError(t, err)
				assert.Equal(t, mockingjay.Endpoints{f.Endpoint}, currentEndpoints)

				for _, request := range f.MatchingRequests {
					t.Run(request.Description, func(t *testing.T) {
						res, report, err := mj.Send(request.Request)
						assert.NoError(t, err)
						assert.True(t, report.HadMatch, pp.Format(report))
						assertResponseMatches(t, f.Endpoint.Response, res)
					})
				}

				for _, request := range f.NonMatchingRequests {
					t.Run(request.Description, func(t *testing.T) {
						_, report, err := mj.Send(request.Request)
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
	assert.Equal(t, want.Body, got.Body)
	assert.Equal(t, want.Status, want.Status)

	for key, v := range want.Headers {
		for _, value := range v {
			i := slices.Index(got.Headers[key], value)
			assert.NotEqual(t, -1, i, fmt.Sprintf("%q not found in %v", value, got.Headers[key]))
		}
	}
}
