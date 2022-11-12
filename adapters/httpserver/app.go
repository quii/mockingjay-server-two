package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

const (
	HeaderMockingjayMatched = "x-mockingjay-matched"
	ReportsPath             = "/match-reports"
	ConfigEndpointsPath     = "/endpoints"
)

type App struct {
	AdminRouter *http.ServeMux

	endpoints    mockingjay.Endpoints
	matchReports []matching.Report
}

func New(endpoints mockingjay.Endpoints) *App {
	app := &App{endpoints: endpoints}

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc(ReportsPath, app.matchReportsHandler)
	adminRouter.HandleFunc(ConfigEndpointsPath, app.configEndpointsHandler)

	app.AdminRouter = adminRouter
	return app
}

func (a *App) StubHandler(w http.ResponseWriter, r *http.Request) {
	matchReport := matching.NewReport(r, a.endpoints)
	a.matchReports = append(a.matchReports, matchReport)

	if !matchReport.HadMatch {
		w.Header().Add(HeaderMockingjayMatched, "false")
		w.Header().Add("content-type", "application/json")
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

func (a *App) matchReportsHandler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(a.matchReports)
}

func (a *App) configEndpointsHandler(w http.ResponseWriter, r *http.Request) {
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
