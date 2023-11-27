package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *Application) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(a.AppName))
	})

	return router
}
