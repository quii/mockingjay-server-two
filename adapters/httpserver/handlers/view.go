package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"cuelang.org/go/cue/cuecontext"
)

const (
	contentTypeJSON = "application/json"
	contentTypeCue  = "application/cue"
	templateDir     = "./adapters/httpserver/handlers"
	templatePattern = "templates/*.gohtml"
)

type View struct {
	templFS fs.FS
	templ   *template.Template
}

func NewContentNegotiatingRenderer(devMode bool) (*View, error) {
	view := View{}
	if devMode {
		view.templFS = os.DirFS(templateDir)
	} else {
		view.templFS = templates
		templ, err := template.ParseFS(view.templFS, templatePattern)
		if err != nil {
			return nil, err
		}
		view.templ = templ
	}
	return &view, nil
}

func (v *View) Render(w http.ResponseWriter, accept string, template string, thing any) {
	switch accept {
	case contentTypeJSON:
		writeJSON(w, thing)
	case contentTypeCue:
		writeCue(w, thing)
	default:
		t, err := v.getTemplates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = t.ExecuteTemplate(w, template, thing)
	}
}

func writeCue(w http.ResponseWriter, content any) {
	c := cuecontext.New()
	w.Header().Add("content-type", contentTypeCue)
	encode := c.Encode(content)
	_, _ = fmt.Fprint(w, encode)
}

func writeJSON(w http.ResponseWriter, content any) {
	w.Header().Add("content-type", contentTypeJSON)
	_ = json.NewEncoder(w).Encode(content)
}

func (v *View) getTemplates() (*template.Template, error) {
	if v.templ != nil {
		return v.templ, nil
	}
	return template.ParseFS(v.templFS, templatePattern)
}
