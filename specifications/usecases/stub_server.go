package usecases

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"golang.org/x/exp/slices"
)

type StubServer struct {
	Admin  Admin
	Client StubServerClient
}

func (s StubServer) Test(t *testing.T, endpoint http.Endpoint) {
	t.Run("stub server for "+endpoint.Description, func(t *testing.T) {
		t.Cleanup(s.mustDeleteEndpoint(t, s.addEndpoint(t, endpoint)))
		report := s.assertEndpointRespondsCorrectly(t, endpoint)
		s.assertReportCanBeFoundFor(t, report.ID)
	})
}

func (s StubServer) assertEndpointRespondsCorrectly(t *testing.T, endpoint http.Endpoint) matching.Report {
	var theReport matching.Report
	t.Run("the endpoint responds correctly to the request it was configured with", func(t *testing.T) {
		res, report, err := s.Client.Send(endpoint.Request)
		assert.True(t, report.HadMatch, report)
		assert.NoError(t, err)
		AssertResponseMatches(t, endpoint.Response, res)
		theReport.ID = report.ID
	})

	return theReport
}

func (s StubServer) addEndpoint(t *testing.T, endpoint http.Endpoint) uuid.UUID {
	var id uuid.UUID
	t.Run("an endpoint can be added", func(t *testing.T) {
		assert.NoError(t, s.Admin.AddEndpoints(endpoint))
		endpoints, err := s.Admin.GetEndpoints()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(endpoints))
		id = endpoints[0].ID
	})
	return id
}

func (s StubServer) assertReportCanBeFoundFor(t *testing.T, id uuid.UUID) {
	t.Run("a report can be found", func(t *testing.T) {
		reports, err := s.Admin.GetReports()
		assert.NoError(t, err)
		i := slices.IndexFunc(reports, func(r matching.Report) bool {
			return r.ID == id
		})
		t.Log(id)
		assert.NotEqual(t, -1, i, pp.Format(reports))
	})
}

func (s StubServer) mustDeleteEndpoint(t *testing.T, id uuid.UUID) func() {
	return func() {
		t.Helper()
		assert.NoError(t, s.Admin.DeleteEndpoint(id))
	}
}
