package usecases

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type StubServerRequestMatching struct {
	Admin  Admin
	Client Client
}

func (s StubServerRequestMatching) Test(t *testing.T, fixture mockingjay.TestFixture) {
	t.Run(fixture.Endpoint.Description, func(t *testing.T) {
		t.Cleanup(s.mustDeleteEndpoint(t, s.mustConfigureEndpoint(t, fixture.Endpoint)))

		for _, request := range fixture.MatchingRequests {
			t.Run("matches on "+request.Description, func(t *testing.T) {
				res, report := s.mustSend(t, request.Request)
				assert.True(t, report.HadMatch, pp.Format(report))
				AssertResponseMatches(t, fixture.Endpoint.Response, res)
			})
		}

		for _, request := range fixture.NonMatchingRequests {
			t.Run("wont match for "+request.Description, func(t *testing.T) {
				_, report := s.mustSend(t, request.Request)
				assert.False(t, report.HadMatch)
			})
		}
	})
}

func (s StubServerRequestMatching) mustConfigureEndpoint(t *testing.T, endpoint http.Endpoint) uuid.UUID {
	t.Helper()
	assert.NoError(t, s.Admin.AddEndpoints(endpoint))
	currentEndpoints, err := s.Admin.GetEndpoints()
	assert.NoError(t, err)
	AssertEndpointsEqual(t, http.Endpoints{endpoint}, currentEndpoints)
	return currentEndpoints[0].ID
}

func (s StubServerRequestMatching) mustSend(t *testing.T, request http.Request) (http.Response, matching.Report) {
	t.Helper()
	res, report, err := s.Client.Send(request)
	assert.NoError(t, err)
	return res, report
}

func (s StubServerRequestMatching) mustDeleteEndpoint(t *testing.T, id uuid.UUID) func() {
	return func() {
		t.Helper()
		assert.NoError(t, s.Admin.DeleteEndpoint(id))
	}
}
