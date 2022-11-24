package handlers

import (
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type StubServerService interface {
	CreateMatchReport(r *http.Request) (matching.Report, error)
}

type StubHandler struct {
	adminBaseURL string
	service      StubServerService
}

func NewStubHandler(service StubServerService, adminBaseURL string) *StubHandler {
	return &StubHandler{adminBaseURL: adminBaseURL, service: service}
}

func (s *StubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchReport, err := s.service.CreateMatchReport(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add(HeaderMockingjayMatchID, matchReport.ID.String())

	if !matchReport.HadMatch {
		w.Header().Add(HeaderMockingjayMatched, "false")
		w.Header().Add("location", s.adminBaseURL+ReportsPath+"/"+matchReport.ID.String())
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeMJResponse(w, matchReport.SuccessfulMatch)
}

func writeMJResponse(w http.ResponseWriter, res stub.Response) {
	w.Header().Add(HeaderMockingjayMatched, "true")

	for key, v := range res.Headers {
		for _, value := range v {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.Status)
	_, _ = w.Write([]byte(res.Body))
}
