package handlers

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"

	http2 "github.com/quii/mockingjay-server-two/domain/mockingjay/http"
	"github.com/quii/mockingjay-server-two/domain/mockingjay/matching"
)

type View struct {
	templFS fs.FS
	templ   *template.Template
}

func (v *View) Endpoints(w http.ResponseWriter, accept string, endpoints []http2.Endpoint) {
	v.render(w, accept, "endpoints.gohtml", endpoints)
}

func (v *View) Reports(w http.ResponseWriter, accept string, reports []matching.Report) {
	v.render(w, accept, "reports.gohtml", reports)
}

func (v *View) Report(w http.ResponseWriter, accept string, report matching.Report) {
	v.render(w, accept, "report.gohtml", report)
}

func (v *View) render(w http.ResponseWriter, accept string, template string, thing any) {
	switch accept {
	case contentTypeApplicationJSON:
		writeJSON(w, thing)
	default:
		t, err := v.getTemplates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = t.ExecuteTemplate(w, template, thing)
	}
}

func writeJSON(w http.ResponseWriter, content any) {
	w.Header().Add("content-type", contentTypeApplicationJSON)
	_ = json.NewEncoder(w).Encode(content)
}

func (v *View) getTemplates() (*template.Template, error) {
	if v.templ != nil {
		return v.templ, nil
	}
	return template.ParseFS(v.templFS, "templates/*.gohtml")
}
