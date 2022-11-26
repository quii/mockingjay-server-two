package usecases

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/cue-exp/cueconfig"
	mj "github.com/quii/mockingjay-server-two"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type (
	RequestDescription struct {
		Description string       `json:"description,omitempty"`
		Request     stub.Request `json:"request"`
	}

	TestFixture struct {
		Endpoint            stub.Endpoint        `json:"endpoint"`
		MatchingRequests    []RequestDescription `json:"matchingRequests,omitempty"`
		NonMatchingRequests []RequestDescription `json:"nonMatchingRequests,omitempty"`
	}
)

type testFixtureCue struct {
	Fixtures []TestFixture
}

func NewFixturesFromCue(basePath string) ([]TestFixture, error) {
	var allFixtures []TestFixture

	dir, err := fs.ReadDir(os.DirFS(basePath), ".")
	if err != nil {
		return nil, err
	}

	for _, f := range dir {
		var fixtures testFixtureCue
		path := basePath + f.Name()
		if err := cueconfig.Load(path, mj.Schema, nil, nil, &fixtures); err != nil {
			return nil, fmt.Errorf("failed to parse %s, %v", path, err)
		}
		allFixtures = append(allFixtures, fixtures.Fixtures...)
	}

	return allFixtures, nil
}
