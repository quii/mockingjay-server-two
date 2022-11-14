package specifications

import (
	"strings"
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

func MockingjayAdmin(t *testing.T, admin Admin, configuration mockingjay.Endpoints) {
	t.Run("can check all endpoints are configured", func(t *testing.T) {
		assert.NoError(t, admin.Configure(configuration...))

		retrievedConfiguration, err := admin.GetCurrentConfiguration()
		assert.NoError(t, err)
		assert.Equal(t, len(configuration), len(retrievedConfiguration))

		removeWhitespaceFromBodies(configuration)
		removeWhitespaceFromBodies(retrievedConfiguration)
		assert.Equal(t, configuration, retrievedConfiguration)
	})
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
