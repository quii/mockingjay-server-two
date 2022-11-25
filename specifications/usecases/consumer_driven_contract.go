package usecases

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type ConsumerDrivenContract struct {
	Admin  Admin
	Client ConsumerDrivenContractChecker
}

func (s ConsumerDrivenContract) Test(t *testing.T, endpoint stub.Endpoint) {
	if len(endpoint.CDCs) == 0 {
		return
	}
	t.Run("cdc for "+endpoint.Description, func(t *testing.T) {
		t.Cleanup(s.mustDeleteEndpoint(t, s.addEndpoint(t, endpoint)))
		results, err := s.Client.CheckEndpoints()
		assert.NoError(t, err)
		for _, result := range results {
			assert.True(t, result.Passed || result.Ignore, pp.Format(results))
		}
	})
}

func (s ConsumerDrivenContract) addEndpoint(t *testing.T, endpoint stub.Endpoint) uuid.UUID {
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

func (s ConsumerDrivenContract) mustDeleteEndpoint(t *testing.T, id uuid.UUID) func() {
	return func() {
		t.Helper()
		assert.NoError(t, s.Admin.DeleteEndpoint(id))
	}
}
