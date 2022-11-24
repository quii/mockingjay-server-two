package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
)

func MockingjayStubServerSpec(
	t *testing.T,
	admin usecases.Admin,
	client usecases.StubServerClient,
	examples stub.Endpoints,
	testFixtures []mockingjay.TestFixture,
) {
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

	for _, f := range testFixtures {
		requestMatchingUseCase.Test(t, f)
	}
}
