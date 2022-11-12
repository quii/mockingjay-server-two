package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
	"golang.org/x/exp/slices"
)

const (
	HeaderMockingjayMatched = "x-mockingjay-matched"
	ReportsPath             = "/match-reports"
	ConfigEndpointsPath     = "/endpoints"
)

type App struct {
	AdminRouter http.Handler

	endpoints    mockingjay.Endpoints
	matchReports map[uuid.UUID]matching.Report
}

func New(endpoints mockingjay.Endpoints) *App {
	app := &App{
		endpoints:    endpoints,
		matchReports: make(map[uuid.UUID]matching.Report),
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc(ReportsPath, app.allReports)
	adminRouter.HandleFunc(fmt.Sprintf("%s/{reportID}", ReportsPath), app.viewReport)
	adminRouter.HandleFunc(ConfigEndpointsPath, app.putEndpoints).Methods(http.MethodPut)
	adminRouter.HandleFunc(ConfigEndpointsPath, app.getEndpoints).Methods(http.MethodGet)

	app.AdminRouter = adminRouter
	return app
}

func (a *App) StubHandler(w http.ResponseWriter, r *http.Request) {
	matchReport := matching.NewReport(r, a.endpoints)

	reportID := uuid.New()
	a.matchReports[reportID] = matchReport

	if !matchReport.HadMatch {
		w.Header().Add(HeaderMockingjayMatched, "false")
		w.Header().Add("location", ReportsPath+"/"+reportID.String())
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
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

func (a *App) allReports(w http.ResponseWriter, r *http.Request) {
	var reports []matching.Report
	for _, report := range a.matchReports {
		reports = append(reports, report)
	}
	slices.SortFunc(reports, func(a, b matching.Report) bool {
		return a.CreatedAt.Before(b.CreatedAt)
	})
	_ = json.NewEncoder(w).Encode(reports)
}

func (a *App) putEndpoints(w http.ResponseWriter, r *http.Request) {
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

func (a *App) viewReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID := vars["reportID"]

	if report, exists := a.matchReports[uuid.MustParse(reportID)]; exists {
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(report)
	} else {
		http.NotFound(w, r)
	}
}

func (a *App) getEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(a.endpoints)
}
