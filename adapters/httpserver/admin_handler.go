package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/quii/mockingjay-server-two/domain/mockingjay"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

const (
	HeaderMockingjayMatched = "x-mockingjay-matched"
	ReportsPath             = "/match-reports"
	ConfigEndpointsPath     = "/endpoints"
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
}

func NewAdminHandler(service AdminServiceService) *AdminHandler {
	app := &AdminHandler{
		service: service,
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc(ReportsPath, app.allReports)
	adminRouter.HandleFunc(fmt.Sprintf("%s/{reportID}", ReportsPath), app.viewReport)
	adminRouter.HandleFunc(ConfigEndpointsPath, app.putEndpoints).Methods(http.MethodPut)
	adminRouter.HandleFunc(ConfigEndpointsPath, app.getEndpoints).Methods(http.MethodGet)

	app.Handler = adminRouter
	return app
}

func (a *AdminHandler) allReports(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(a.service.GetReports())
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

func (a *AdminHandler) viewReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := uuid.Parse(mux.Vars(r)["reportID"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if report, exists := a.service.GetReport(reportID); exists {
		w.Header().Add("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(report)
	} else {
		http.NotFound(w, r)
	}
}

func (a *AdminHandler) getEndpoints(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(a.service.GetEndpoints())
}
