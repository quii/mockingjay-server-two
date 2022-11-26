package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
)

const (
	examplesDir = "/examples/"
	fixturesDir = "/stub_server_fixtures/"
)

func MockingjayStubServerSpec(
	t *testing.T,
	admin usecases.Admin,
	client usecases.StubServerClient,
	specRoot string,
) {
	fixtures, err := usecases.NewFixturesFromCue(specRoot + fixturesDir)
	assert.NoError(t, err)
	examples, err := stub.NewEndpointsFromCue(specRoot + examplesDir)
	assert.NoError(t, err)

	assert.NoError(t, admin.DeleteEndpoints())
	assert.NoError(t, admin.DeleteReports())

	stubServerUseCase := usecases.StubServer{
		Admin:  admin,
		Client: client,
	}
	requestMatchingUseCase := usecases.StubServerRequestMatching{
		Admin:  admin,
		Client: client,
	}

	for _, example := range examples {
		stubServerUseCase.Test(t, example)
	}

	for _, f := range fixtures {
		requestMatchingUseCase.Test(t, f)
	}
}
