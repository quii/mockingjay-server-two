package handlers

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

type View struct {
	templFS fs.FS
	templ   *template.Template
}

func NewContentNegotiatingRenderer(devMode bool) (*View, error) {
	view := View{}
	if devMode {
		view.templFS = os.DirFS("./adapters/httpserver/handlers")
	} else {
		view.templFS = templates
		templ, err := template.ParseFS(view.templFS, "templates/*.gohtml")
		if err != nil {
			return nil, err
		}
		view.templ = templ
	}
	return &view, nil
}

func (v *View) Render(w http.ResponseWriter, accept string, template string, thing any) {
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
