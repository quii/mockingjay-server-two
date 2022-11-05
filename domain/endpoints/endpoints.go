package endpoints

import (
	"errors"
	"net/http"

	"github.com/cue-exp/cueconfig"
	mj "github.com/quii/mockingjay-server-two"
	"golang.org/x/exp/slices"
)

var (
	ErrNoMatchingRequests = errors.New("no matching requests")
)

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

func NewEndpointsFromCue(config string) (Endpoints, error) {
	var endpoints Endpoints

	if err := cueconfig.Load(config, mj.Schema, nil, nil, &endpoints); err != nil {
		return Endpoints{}, err
	}

	return endpoints, nil
}

func (e Endpoints) FindMatchingResponse(req *http.Request) (Response, error) {
	i := slices.IndexFunc(e.Endpoints, EndpointMatcher(req))
	if i == -1 {
		return Response{}, ErrNoMatchingRequests
	}
	return e.Endpoints[i].Response, nil
}

func EndpointMatcher(req *http.Request) func(Endpoint) bool {
	return func(e Endpoint) bool {
		methodMatch := e.Request.Method == req.Method
		pathMatch := e.Request.Path == req.URL.Path

		return methodMatch && pathMatch
	}
}
