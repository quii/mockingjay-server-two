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
		assert.Equal(t, examples[0].Description, configuration[0].Description)
		assert.Equal(t, examples[0].Request.Method, configuration[0].Request.Method)
		assert.Equal(t, examples[0].Request.RegexPath, configuration[0].Request.RegexPath)
		assert.Equal(t, examples[0].Request.Path, configuration[0].Request.Path)
		assert.Equal(t, examples[0].Request.Body, configuration[0].Request.Body)
		assert.Equal(t, examples[0].Response.Status, configuration[0].Response.Status)
		assert.Equal(t, examples[0].Response.Body, configuration[0].Response.Body)
	})
}
