package specifications

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"golang.org/x/exp/slices"
)

type StubServerUseCase struct {
	admin  Admin
	client Client
}

func (s StubServerUseCase) Test(t *testing.T, endpoint http.Endpoint) {
	t.Run("for "+endpoint.Description, func(t *testing.T) {
		id := s.addEndpoint(t, endpoint)

		t.Cleanup(func() {
			assert.NoError(t, s.admin.DeleteEndpoint(id))
		})

		report := s.assertEndpointRespondsCorrectly(t, endpoint)
		s.assertReportCanBeFoundFor(t, report.ID)
	})
}

func (s StubServerUseCase) assertEndpointRespondsCorrectly(t *testing.T, endpoint http.Endpoint) matching.Report {
	var theReport matching.Report
	t.Run("the endpoint responds correctly to the request it was configured with", func(t *testing.T) {
		res, report, err := s.client.Send(endpoint.Request)
		assert.True(t, report.HadMatch, report)
		assert.NoError(t, err)
		AssertResponseMatches(t, endpoint.Response, res)
		theReport.ID = report.ID
	})

	return theReport
}

func (s StubServerUseCase) addEndpoint(t *testing.T, endpoint http.Endpoint) uuid.UUID {
	var id uuid.UUID
	t.Run("an endpoint can be added", func(t *testing.T) {
		assert.NoError(t, s.admin.AddEndpoints(endpoint))
		endpoints, err := s.admin.GetEndpoints()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(endpoints))
		id = endpoints[0].ID
	})
	return id
}

func (s StubServerUseCase) assertReportCanBeFoundFor(t *testing.T, id uuid.UUID) {
	t.Run("a report can be found", func(t *testing.T) {
		reports, err := s.admin.GetReports()
		assert.NoError(t, err)
		i := slices.IndexFunc(reports, func(r matching.Report) bool {
			return r.ID == id
		})
		t.Log(id)
		assert.NotEqual(t, -1, i, pp.Format(reports))
	})
}
