package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

type EndpointMatcher interface {
	GetMatchReport(r *http.Request) endpoints.MatchReport
}

type StubServer struct {
	matcher EndpointMatcher
}

func NewStubServer(matcher EndpointMatcher) *StubServer {
	return &StubServer{matcher: matcher}
}

func (s StubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchReport := s.matcher.GetMatchReport(r)
	res, exists := matchReport.FindMatchingResponse()

	if !exists {
		if err := json.NewEncoder(w).Encode(matchReport); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for key, v := range res.Headers {
		for _, value := range v {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.Status)
	_, _ = w.Write([]byte(res.Body))
}
