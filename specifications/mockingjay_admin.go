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

func MockingjayAdmin(t *testing.T, admin Admin, examples mockingjay.Endpoints) {
	t.Run("can check all endpoints are configured", func(t *testing.T) {
		assert.NoError(t, admin.Configure(examples...))

		configuration, err := admin.GetCurrentConfiguration()
		assert.NoError(t, err)
		assert.Equal(t, len(examples), len(configuration))

		for i := range examples {
			assert.Equal(t, examples[i].Description, configuration[i].Description)

			originalRequest := examples[i].Request
			retrievedRequest := configuration[i].Request
			assert.Equal(t, originalRequest.Method, retrievedRequest.Method)
			assert.Equal(t, originalRequest.RegexPath, retrievedRequest.RegexPath)
			assert.Equal(t, originalRequest.Path, retrievedRequest.Path)
			assert.Equal(t, originalRequest.Body, retrievedRequest.Body)
			assert.Equal(t, originalRequest.Headers, retrievedRequest.Headers)

			originalResponse := examples[i].Response
			retrievedResponse := configuration[i].Response
			assert.Equal(t, originalResponse.Status, retrievedResponse.Status)
			assert.Equal(t, fudgeTheWhiteSpace(originalResponse.Body), fudgeTheWhiteSpace(retrievedResponse.Body))
			assert.Equal(t, originalResponse.Headers, retrievedResponse.Headers)
		}
	})
}

func fudgeTheWhiteSpace(in string) string {
	in = strings.Replace(in, "\t", "", -1)
	in = strings.Replace(in, "\n", "", -1)
	in = strings.Replace(in, " ", "", -1)
	return in
}
