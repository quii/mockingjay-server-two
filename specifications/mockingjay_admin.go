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

		for i := range examples {
			assert.Equal(t, examples[i].Description, configuration[i].Description)
			assert.Equal(t, examples[i].Request.Method, configuration[i].Request.Method)
			assert.Equal(t, examples[i].Request.RegexPath, configuration[i].Request.RegexPath)
			assert.Equal(t, examples[i].Request.Path, configuration[i].Request.Path)
			assert.Equal(t, examples[i].Request.Body, configuration[i].Request.Body)
			assert.Equal(t, examples[i].Request.Headers, configuration[i].Request.Headers)

			assert.Equal(t, examples[i].Response.Status, configuration[i].Response.Status)
			//assert.Equal(t, examples[i].Response.Body, configuration[i].Response.Body)
			assert.Equal(t, examples[i].Response.Headers, configuration[i].Response.Headers)
		}

	})
}
