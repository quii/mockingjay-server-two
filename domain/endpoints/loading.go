package endpoints

import (
	"fmt"
	"io/fs"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
)

type endpointsCue struct {
	Endpoints []Endpoint
}

func NewEndpointsFromCue(basePath string, configDir fs.FS) (Endpoints, error) {
	var allEndpoints endpointsCue

	dir, err := fs.ReadDir(configDir, ".")
	if err != nil {
		return Endpoints{}, err
	}

	for _, f := range dir {
		var endpoints endpointsCue
		path := basePath + f.Name()
		if err := cueconfig.Load(path, mj.Schema, nil, nil, &endpoints); err != nil {
			return Endpoints{}, fmt.Errorf("failed to parse %s, %v", path, err)
		}
		allEndpoints.Endpoints = append(allEndpoints.Endpoints, endpoints.Endpoints...)
	}

	return allEndpoints.Endpoints, nil
}
