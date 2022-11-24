package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/stub"
)

type EndpointHandler struct {
	service  AdminServiceService
	renderer HTTPRenderer
}

func (a *EndpointHandler) allEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := a.service.Endpoints().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.renderer.Render(w, r.Header.Get("Accept"), "endpoints.gohtml", endpoints)
}

func (a *EndpointHandler) addEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == contentTypeApplicationJSON {
		var newEndpoint stub.Endpoint
		if err := json.NewDecoder(r.Body).Decode(&newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := newEndpoint.Compile(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newEndpoint.LoadedAt = time.Now().UTC()

		if err := a.service.Endpoints().Create(newEndpoint.ID, newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestHeaders := make(stub.Headers)
		if r.FormValue("request.header.name") != "" {
			requestHeaders[r.FormValue("request.header.name")] = strings.Split(r.FormValue("request.header.values"), "; ")
		}

		responseHeaders := make(stub.Headers)
		if r.FormValue("response.header.name") != "" {
			responseHeaders[r.FormValue("response.header.name")] = strings.Split(r.FormValue("response.header.values"), "; ")
		}

		status, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newEndpoint := stub.Endpoint{
			ID:          uuid.New(),
			Description: r.FormValue("description"),
			Request: stub.Request{
				Method:    r.FormValue("method"),
				RegexPath: r.FormValue("regexpath"),
				Path:      r.FormValue("path"),
				Headers:   requestHeaders,
				Body:      r.FormValue("request.body"),
			},
			Response: stub.Response{
				Status:  status,
				Body:    r.FormValue("response.body"),
				Headers: responseHeaders,
			},
			LoadedAt: time.Now().UTC(),
			CDCs:     nil,
		}

		if err := newEndpoint.Compile(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.service.Endpoints().Create(newEndpoint.ID, newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func (a *EndpointHandler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["endpointIndex"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.service.Endpoints().Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
