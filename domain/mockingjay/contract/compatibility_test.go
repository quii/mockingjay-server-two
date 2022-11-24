package contract_test

import (
	"testing"

	"github.com/adamluzsi/testcase/pp"
	"github.com/alecthomas/assert/v2"
	"github.com/cue-exp/cueconfig"
	mj "github.com/quii/mockingjay-server-two"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/contract"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type CompatibilityFixture struct {
	Description        string
	Got, Want          stub.Response
	ShouldBeCompatible bool
}

type fixturesCue struct {
	CDCFixtures []CompatibilityFixture
}

func TestIsResponseCompatible(t *testing.T) {
	var fixtures fixturesCue
	assert.NoError(t, cueconfig.Load("compatibility_tests.cue", mj.Schema, nil, nil, &fixtures))

	for _, fixture := range fixtures.CDCFixtures {
		t.Run(fixture.Description, func(t *testing.T) {
			assert.Equal(t, fixture.ShouldBeCompatible, contract.IsResponseCompatible(fixture.Got, fixture.Want), pp.Diff(fixture.Got, fixture.Want))
		})
	}
}
