package handlers

import (
	"net/http"
)

type CDCHandler struct {
	service  AdminServiceService
	renderer HTTPRenderer
}

func (a *CDCHandler) checkContracts(w http.ResponseWriter, r *http.Request) {
	reports, err := a.service.CheckEndpoints()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.renderer.Render(w, r.Header.Get("Accept"), "cdc_reports.gohtml", reports)
}
