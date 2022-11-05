package endpoints

import (
	"io/fs"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
)

func NewEndpointsFromCue(basePath string, configDir fs.FS) (Endpoints, error) {
	var allEndpoints Endpoints

	dir, err := fs.ReadDir(configDir, ".")
	if err != nil {
		return Endpoints{}, err
	}

	for _, f := range dir {
		var endpoints Endpoints
		if err := cueconfig.Load(basePath+f.Name(), mj.Schema, nil, nil, &endpoints); err != nil {
			return Endpoints{}, err
		}
		allEndpoints.Endpoints = append(allEndpoints.Endpoints, endpoints.Endpoints...)
	}

	return allEndpoints, nil
}
