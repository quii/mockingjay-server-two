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

func (a *App) StubHandler(w http.ResponseWriter, r *http.Request) {
	matchReport := matching.NewReport(r, a.endpoints)

	if !matchReport.HadMatch {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(matchReport); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

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
	if err := json.NewDecoder(r.Body).Decode(&a.endpoints); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
