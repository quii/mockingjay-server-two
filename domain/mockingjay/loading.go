package mockingjay

import (
	"fmt"
	"io/fs"

	"github.com/cue-exp/cueconfig"
	"github.com/quii/mockingjay-server-two"
)

type endpointsCue struct {
	Endpoints []Endpoint
}

type testFixtureCue struct {
	Fixtures []TestFixture
}

//func load[DTO any, Data any](basePath string, configDir fs.FS) (Data, error) {
//	var allData Data
//
//	dir, err := fs.ReadDir(configDir, ".")
//	if err != nil {
//		return allData, err
//	}
//
//	for _, f := range dir {
//		var data DTO
//		path := basePath + f.Name()
//		if err := cueconfig.Load(path, mj.Schema, nil, nil, &data); err != nil {
//			return allData, fmt.Errorf("failed to parse %s, %v", path, err)
//		}
//		allData.Endpoints = append(allEndpoints.Endpoints, endpoints.Endpoints...)
//	}
//}

func NewEndpointsFromCue(basePath string, configDir fs.FS) (Endpoints, error) {
	var allEndpoints Endpoints

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
		allEndpoints = append(allEndpoints, endpoints.Endpoints...)
	}

	return allEndpoints, nil
}

func NewFixturesFromCue(basePath string, configDir fs.FS) ([]TestFixture, error) {
	var allFixtures []TestFixture

	dir, err := fs.ReadDir(configDir, ".")
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
