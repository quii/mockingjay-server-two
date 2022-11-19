package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Admin interface {
	AddEndpoints(endpoints ...mockingjay.Endpoint) error
	GetReports() ([]matching.Report, error)
	Reset() error
	GetEndpoints() (mockingjay.Endpoints, error)
	DeleteEndpoint(uuid uuid.UUID) error
}

func MockingjayAdmin(t *testing.T, admin Admin, endpoints mockingjay.Endpoints) {
	t.Run("can configure new endpoints and retrieve them", func(t *testing.T) {
		assert.NoError(t, admin.Reset())
		assert.NoError(t, admin.AddEndpoints(endpoints...))

		retrievedEndpoints, err := admin.GetEndpoints()
		assert.NoError(t, err)
		AssertEndpointsEqual(t, retrievedEndpoints, endpoints)

		t.Run("and can then delete them", func(t *testing.T) {
			for _, endpoint := range retrievedEndpoints {
				assert.NoError(t, admin.DeleteEndpoint(endpoint.ID))
			}
			gotEndpoints, err := admin.GetEndpoints()

			assert.NoError(t, err)
			assert.Equal(t, 0, len(gotEndpoints))
		})
	})
}
