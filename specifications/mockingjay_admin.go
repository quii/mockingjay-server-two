package specifications

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/google/uuid"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Admin interface {
	Configure(endpoints ...mockingjay.Endpoint) error
	GetReports() ([]matching.Report, error)
	Reset() error
	GetEndpoints() (mockingjay.Endpoints, error)
}

func MockingjayAdmin(t *testing.T, admin Admin, endpoints mockingjay.Endpoints) {
	t.Run("can check all endpoints are configured", func(t *testing.T) {
		assert.NoError(t, admin.Reset())
		assert.NoError(t, admin.Configure(endpoints...))

		retrievedConfiguration, err := admin.GetEndpoints()
		assert.NoError(t, err)
		assert.Equal(t, len(endpoints), len(retrievedConfiguration))

		//todo: make an AssertEqual method on endpoints or something to tidy this up
		removeWhitespaceFromBodies(endpoints)
		removeWhitespaceFromBodies(retrievedConfiguration)
		zeroUUIDs(endpoints)
		zeroUUIDs(retrievedConfiguration)

		//todo: to test all endpoints will require support for multiple headers
		for i, endpoint := range endpoints {
			assert.Equal(t, endpoint, retrievedConfiguration[i])
		}
	})
}

func zeroUUIDs(endpoints mockingjay.Endpoints) {
	for i := range endpoints {
		endpoints[i].ID = uuid.UUID{}
	}
}

func removeWhitespaceFromBodies(endpoints mockingjay.Endpoints) {
	for i := range endpoints {
		endpoints[i].Response.Body = fudgeTheWhiteSpace(endpoints[i].Response.Body)
	}
}

func fudgeTheWhiteSpace(in string) string {
	in = strings.Replace(in, "\t", "", -1)
	in = strings.Replace(in, "\n", "", -1)
	in = strings.Replace(in, " ", "", -1)
	return in
}
