package stub

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
)

type (
	EndpointsCue struct {
		Endpoints []Endpoint
	}
)

func NewEndpointsFromCue(basePath string) (Endpoints, error) {
	var allEndpoints Endpoints

	dir, err := fs.ReadDir(os.DirFS(basePath), ".")
	if err != nil {
		return Endpoints{}, err
	}

	for _, f := range dir {
		var endpoints EndpointsCue
		path := basePath + f.Name()
		if err := cueconfig.Load(path, mj.Schema, nil, nil, &endpoints); err != nil {
			return Endpoints{}, fmt.Errorf("failed to parse %s, %v", path, err)
		}
		allEndpoints = append(allEndpoints, endpoints.Endpoints...)
	}

	return allEndpoints, nil
}
