package handlers

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

var (
	//go:embed "templates/*"
	templates embed.FS
)

const (
	HeaderMockingjayMatched = "x-mockingjay-matched"
	ReportsPath             = "/match-reports"
	EndpointsPath           = "/"
)

type AdminServiceService interface {
	GetReports() matching.Reports
	GetReport(id uuid.UUID) (matching.Report, bool)
	PutEndpoints(e mockingjay.Endpoints) error
	GetEndpoints() mockingjay.Endpoints
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
	adminRouter.HandleFunc(EndpointsPath, app.putEndpoints).Methods(http.MethodPut)
	adminRouter.HandleFunc(EndpointsPath, app.getEndpoints).Methods(http.MethodGet)
	adminRouter.HandleFunc(EndpointsPath, app.deleteEndpoints).Methods(http.MethodDelete)
	adminRouter.HandleFunc(EndpointsPath, app.addEndpoint).Methods(http.MethodPost)

	app.Handler = adminRouter
	return app
}

func (a *AdminHandler) allReports(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, a.service.GetReports())
}

func (a *AdminHandler) getEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints := a.service.GetEndpoints()
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

	if report, exists := a.service.GetReport(reportID); exists {
		writeJSON(w, report)
		return
	}

	http.NotFound(w, r)
}

func (a *AdminHandler) putEndpoints(w http.ResponseWriter, r *http.Request) {
	var newEndpoints mockingjay.Endpoints
	if err := json.NewDecoder(r.Body).Decode(&newEndpoints); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.service.PutEndpoints(newEndpoints); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *AdminHandler) addEndpoint(w http.ResponseWriter, r *http.Request) {
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
		requestHeaders[r.FormValue("response.header.name")] = strings.Split(r.FormValue("response.header.values"), "; ")
	}

	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newEndpoint := mockingjay.Endpoint{
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

	if err := a.service.PutEndpoints(append(a.service.GetEndpoints(), newEndpoint)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (a *AdminHandler) deleteEndpoints(w http.ResponseWriter, _ *http.Request) {
	if err := a.service.PutEndpoints(nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, content any) {
	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(content)
}
