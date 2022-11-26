package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
	"github.com/quii/mockingjay-server-two/specifications/usecases"
)

func MockingjayConsumerDrivenContractSpec(t *testing.T, admin usecases.Admin, cdcClient usecases.ConsumerDrivenContractChecker, specroot string) {
	examples, err := stub.NewEndpointsFromCue(specroot + examplesDir)
	assert.NoError(t, err)

	assert.NoError(t, admin.DeleteEndpoints())
	assert.NoError(t, admin.DeleteReports())

	cdcUseCase := usecases.ConsumerDrivenContract{
		Admin:  admin,
		Client: cdcClient,
	}

	for _, example := range examples {
		cdcUseCase.Test(t, example)
	}
}
