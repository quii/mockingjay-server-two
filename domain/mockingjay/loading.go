package mockingjay

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
)

type endpointsCue struct {
	Endpoints []Endpoint
}

type testFixtureCue struct {
	Fixtures []TestFixture
}

func NewEndpointsFromCue(basePath string) (Endpoints, error) {
	var allEndpoints Endpoints

	dir, err := fs.ReadDir(os.DirFS(basePath), ".")
	if err != nil {
		return Endpoints{}, err
	}

	for _, f := range dir {
		var endpoints endpointsCue
		path := basePath + f.Name()
		if err := cueconfig.Load(path, mj.Schema, nil, nil, &endpoints); err != nil {
			return Endpoints{}, fmt.Errorf("failed to parse %s, %v", path, err)
		}
		allEndpoints = append(allEndpoints, endpoints.Endpoints...)
	}

	return allEndpoints, nil
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
		if err := cueconfig.Load(path, mj.FixtureSchema, nil, nil, &fixtures); err != nil {
			return nil, fmt.Errorf("failed to parse %s, %v", path, err)
		}
		allFixtures = append(allFixtures, fixtures.Fixtures...)
	}

	return allFixtures, nil
}
