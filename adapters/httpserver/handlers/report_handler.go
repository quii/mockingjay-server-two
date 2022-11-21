package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReportHandler struct {
	service  AdminServiceService
	renderer HTTPRenderer
}

func (a *ReportHandler) allReports(w http.ResponseWriter, r *http.Request) {
	reports, err := a.service.Reports().GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.renderer.Render(w, r.Header.Get("Accept"), "reports.gohtml", reports)
}

func (a *ReportHandler) getReport(w http.ResponseWriter, r *http.Request) {
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
	a.renderer.Render(w, r.Header.Get("Accept"), "report.gohtml", report)
}

func (a *ReportHandler) deleteReports(w http.ResponseWriter, _ *http.Request) {
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
