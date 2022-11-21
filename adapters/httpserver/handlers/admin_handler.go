package handlers

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/crud"
	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

var (
	//go:embed "templates/*"
	templates embed.FS

	//go:embed static
	static embed.FS
)

const (
	HeaderMockingjayMatched    = "x-mockingjay-matched"
	ReportsPath                = "/match-reports"
	EndpointsPath              = "/"
	contentTypeApplicationJSON = "application/json"
)

type AdminServiceService interface {
	Reports() crud.CRUDesque[uuid.UUID, matching.Report]
	Endpoints() crud.CRUDesque[uuid.UUID, http2.Endpoint]
}

type AdminHandler struct {
	http.Handler
	service AdminServiceService
	view    View
}

func NewAdminHandler(service AdminServiceService, devMode bool) (*AdminHandler, error) {
	app := &AdminHandler{
		service: service,
		view:    View{},
	}

	if devMode {
		app.view.templFS = os.DirFS("./adapters/httpserver/handlers")
	} else {
		app.view.templFS = templates
		templ, err := template.ParseFS(app.view.templFS, "templates/*.gohtml")
		if err != nil {
			return nil, err
		}
		app.view.templ = templ
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc(ReportsPath, app.allReports).Methods(http.MethodGet)
	adminRouter.HandleFunc(ReportsPath, app.deleteReports).Methods(http.MethodDelete)
	adminRouter.HandleFunc(ReportsPath+"/{reportID}", app.getReport).Methods(http.MethodGet)

	adminRouter.HandleFunc(EndpointsPath, app.getEndpoints).Methods(http.MethodGet)
	adminRouter.HandleFunc(EndpointsPath+"{endpointIndex}", app.DeleteEndpoint).Methods(http.MethodDelete)
	adminRouter.HandleFunc(EndpointsPath, app.addEndpoint).Methods(http.MethodPost)

	lol, _ := fs.Sub(static, "static")
	adminRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(lol))))

	app.Handler = adminRouter
	return app, nil
}

func (a *AdminHandler) getEndpoints(w http.ResponseWriter, r *http.Request) {
	endpoints, err := a.service.Endpoints().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.view.Endpoints(w, r.Header.Get("Accept"), endpoints)
}

func (a *AdminHandler) allReports(w http.ResponseWriter, r *http.Request) {
	reports, err := a.service.Reports().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.view.Reports(w, r.Header.Get("Accept"), reports)
}

func (a *AdminHandler) getReport(w http.ResponseWriter, r *http.Request) {
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
	a.view.Report(w, r.Header.Get("Accept"), report)
}

func (a *AdminHandler) addEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") == contentTypeApplicationJSON {
		var newEndpoint http2.Endpoint
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

		requestHeaders := make(http2.Headers)
		if r.FormValue("request.header.name") != "" {
			requestHeaders[r.FormValue("request.header.name")] = strings.Split(r.FormValue("request.header.values"), "; ")
		}

		responseHeaders := make(http2.Headers)
		if r.FormValue("response.header.name") != "" {
			responseHeaders[r.FormValue("response.header.name")] = strings.Split(r.FormValue("response.header.values"), "; ")
		}

		status, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newEndpoint := http2.Endpoint{
			ID:          uuid.New(),
			Description: r.FormValue("description"),
			Request: http2.Request{
				Method:    r.FormValue("method"),
				RegexPath: r.FormValue("regexpath"),
				Path:      r.FormValue("path"),
				Headers:   requestHeaders,
				Body:      r.FormValue("request.body"),
			},
			Response: http2.Response{
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

func (a *AdminHandler) deleteReports(w http.ResponseWriter, _ *http.Request) {
	reports, err := a.service.Reports().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, report := range reports {
		if err := a.service.Reports().Delete(report.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
