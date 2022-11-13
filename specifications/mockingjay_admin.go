package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type Admin interface {
	Configure(endpoints ...mockingjay.Endpoint) error
	GetReports() ([]matching.Report, error)
	GetCurrentConfiguration() (mockingjay.Endpoints, error)
}

func MockingjayAdmin(t *testing.T, admin Admin, examples mockingjay.Endpoints) {
	t.Run("can check all endpoints are configured", func(t *testing.T) {
		assert.NoError(t, admin.Configure(examples...))

		configuration, err := admin.GetCurrentConfiguration()
		assert.NoError(t, err)
		assert.Equal(t, len(examples), len(configuration))
	})
}
