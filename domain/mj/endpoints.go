package mj

import (
	"errors"
	"net/http"

	"golang.org/x/exp/slices"
)

var (
	ErrNoMatchingRequests = errors.New("no matching requests")
)

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints,omitempty"`
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
