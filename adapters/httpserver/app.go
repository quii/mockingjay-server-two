package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type App struct {
	endpoints mockingjay.Endpoints
}

const HeaderMockingjayMatched = "x-mockingjay-matched"

func (a *App) StubHandler(w http.ResponseWriter, r *http.Request) {
	matchReport := matching.NewReport(r, a.endpoints)

	if !matchReport.HadMatch {
		w.Header().Add(HeaderMockingjayMatched, "false")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(matchReport); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Add(HeaderMockingjayMatched, "true")
	res := matchReport.SuccessfulMatch
	for key, v := range res.Headers {
		for _, value := range v {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.Status)
	_, _ = w.Write([]byte(res.Body))
}

func (a *App) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	var newEndpoints mockingjay.Endpoints
	if err := json.NewDecoder(r.Body).Decode(&newEndpoints); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.endpoints = newEndpoints

	for i := range a.endpoints {
		if err := a.endpoints[i].Request.Compile(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
