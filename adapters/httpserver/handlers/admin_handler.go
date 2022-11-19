package handlers

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/crud"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

var (
	//go:embed "templates/*"
	templates embed.FS

	//go:embed static
	static embed.FS
)

const (
	HeaderMockingjayMatched = "x-mockingjay-matched"
	ReportsPath             = "/match-reports"
	EndpointsPath           = "/"
)

type AdminServiceService interface {
	Reset()
	Reports() crud.CRUD[uuid.UUID, matching.Report]
	Endpoints() crud.CRUD[uuid.UUID, mockingjay.Endpoint]
}

type AdminHandler struct {
	http.Handler
	service AdminServiceService
	templ   *template.Template
}

func NewAdminHandler(service AdminServiceService) *AdminHandler {
	templ, err := template.ParseFS(templates, "templates/*.gohtml")
	if err != nil {
		panic(err) //todo: fixme
	}

	app := &AdminHandler{
		service: service,
		templ:   templ,
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc(ReportsPath, app.allReports).Methods(http.MethodGet)
	adminRouter.HandleFunc(ReportsPath+"/{reportID}", app.viewReport).Methods(http.MethodGet)

	adminRouter.HandleFunc(EndpointsPath, app.getEndpoints).Methods(http.MethodGet)
	adminRouter.HandleFunc(EndpointsPath, app.deleteEndpoints).Methods(http.MethodDelete)
	adminRouter.HandleFunc(EndpointsPath+"{endpointIndex}", app.DeleteEndpoint).Methods(http.MethodDelete)
	adminRouter.HandleFunc(EndpointsPath, app.addEndpoint).Methods(http.MethodPost)

	lol, _ := fs.Sub(static, "static")
	adminRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(lol))))

	app.Handler = adminRouter
	return app
}

func (a *AdminHandler) allReports(w http.ResponseWriter, _ *http.Request) {
	endpoints, err := a.service.Endpoints().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, endpoints)
}

func (a *AdminHandler) getEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := a.service.Endpoints().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Header.Get("Accept") == "application/json" {
		writeJSON(w, endpoints)
	} else {
		_ = a.templ.ExecuteTemplate(w, "endpoints.gohtml", endpoints)
	}
}

func (a *AdminHandler) viewReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := uuid.Parse(mux.Vars(r)["reportID"])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	report, exists, err := a.service.Reports().GetByID(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, report)
}

func (a *AdminHandler) addEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == "application/json" {
		var newEndpoint mockingjay.Endpoint
		if err := json.NewDecoder(r.Body).Decode(&newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := newEndpoint.Request.Compile(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.service.Endpoints().Create(newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestHeaders := make(mockingjay.Headers)
		if r.FormValue("request.header.name") != "" {
			requestHeaders[r.FormValue("request.header.name")] = strings.Split(r.FormValue("request.header.values"), "; ")
		}

		responseHeaders := make(mockingjay.Headers)
		if r.FormValue("response.header.name") != "" {
			responseHeaders[r.FormValue("response.header.name")] = strings.Split(r.FormValue("response.header.values"), "; ")
		}

		status, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newEndpoint := mockingjay.Endpoint{
			ID:          uuid.New(),
			Description: r.FormValue("description"),
			Request: mockingjay.Request{
				Method:    r.FormValue("method"),
				RegexPath: r.FormValue("regexpath"),
				Path:      r.FormValue("path"),
				Headers:   requestHeaders,
				Body:      r.FormValue("request.body"),
			},
			Response: mockingjay.Response{
				Status:  status,
				Body:    r.FormValue("response.body"),
				Headers: responseHeaders,
			},
			CDCs: nil,
		}

		if err := a.service.Endpoints().Create(newEndpoint); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func (a *AdminHandler) deleteEndpoints(w http.ResponseWriter, _ *http.Request) {
	a.service.Reset()
	w.WriteHeader(http.StatusNoContent)
}

func (a *AdminHandler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
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

func writeJSON(w http.ResponseWriter, content any) {
	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(content)
}
