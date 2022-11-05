package httpserver

import (
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

type EndpointMatcher interface {
	FindMatchingResponse(r *http.Request) (endpoints.Response, error)
}

type StubServer struct {
	matcher EndpointMatcher
}

func NewStubServer(matcher EndpointMatcher) *StubServer {
	return &StubServer{matcher: matcher}
}

func (s StubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := s.matcher.FindMatchingResponse(r)

	if err == endpoints.ErrNoMatchingRequests {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(res.Status)
	_, _ = w.Write([]byte(res.Body))
}
