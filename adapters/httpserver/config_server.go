package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/endpoints"
)

func ConfigServer(endpoints *endpoints.Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&endpoints); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
