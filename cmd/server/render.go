package server

import (
	"fmt"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

type TemplateData struct {
	URL             string
	IsAuthenticated bool
	AuthUser        string
	Flash           string
	Error           string
	CSRFToken       string
}

func (a *Application) DefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.URL = a.Server.Url
	return td
}

func (a *Application) Render(w http.ResponseWriter, r *http.Request, view string, v jet.VarMap) error {
	td := &TemplateData{}
	td = a.DefaultData(td, r)

	tp, err := a.View.GetTemplate(fmt.Sprintf("%s.html", view))
	if err != nil {
		return err
	}

	if err = tp.Execute(w, v, td); err != nil {
		return err
	}

	return nil
}
