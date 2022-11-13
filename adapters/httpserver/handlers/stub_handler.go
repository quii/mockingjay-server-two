package handlers

import (
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type StubServerService interface {
	GetMatchReport(r *http.Request) matching.Report
}

type StubHandler struct {
	adminBaseURL string
	service      StubServerService
}

func NewStubHandler(service StubServerService, adminBaseURL string) *StubHandler {
	return &StubHandler{adminBaseURL: adminBaseURL, service: service}
}

func (s *StubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchReport := s.service.GetMatchReport(r)

	if !matchReport.HadMatch {
		w.Header().Add(HeaderMockingjayMatched, "false")
		w.Header().Add("location", s.adminBaseURL+ReportsPath+"/"+matchReport.ID.String())
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeMJResponse(w, matchReport.SuccessfulMatch)
}

func writeMJResponse(w http.ResponseWriter, res mockingjay.Response) {
	w.Header().Add(HeaderMockingjayMatched, "true")

	for key, v := range res.Headers {
		for _, value := range v {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.Status)
	_, _ = w.Write([]byte(res.Body))
}
