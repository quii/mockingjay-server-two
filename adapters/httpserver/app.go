package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

type App struct {
	endpoints endpoints.Endpoints
}

func (a *App) StubHandler(w http.ResponseWriter, r *http.Request) {
	matchReport := a.endpoints.GetMatchReport(r)
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

func (a *App) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&a.endpoints); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}