package mockingjay

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/http"
)

type (
	endpointsCue struct {
		Endpoints []http.Endpoint
	}
	testFixtureCue struct {
		Fixtures []TestFixture
	}
)

func NewEndpointsFromCue(basePath string) (http.Endpoints, error) {
	var allEndpoints http.Endpoints

	dir, err := fs.ReadDir(os.DirFS(basePath), ".")
	if err != nil {
		return http.Endpoints{}, err
	}

	for _, f := range dir {
		var endpoints endpointsCue
		path := basePath + f.Name()
		if err := cueconfig.Load(path, mj.Schema, nil, nil, &endpoints); err != nil {
			return http.Endpoints{}, fmt.Errorf("failed to parse %s, %v", path, err)
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
