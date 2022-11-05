package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/config"
)

func main() {
	var endpoints config.Endpoints
	if err := http.ListenAndServe(":8080", veryBasicMJServer(endpoints)); err != nil {
		log.Fatal(err)
	}
}

func veryBasicMJServer(endpoints config.Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&endpoints); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		} else {
			for _, endpoint := range endpoints.Endpoints {
				if endpoint.Request.Method == r.Method && endpoint.Request.Path == r.URL.Path {
					_, _ = w.Write([]byte(endpoint.Response.Body))
					w.WriteHeader(endpoint.Response.Status)
				}
			}
		}
	}
}
