package server

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Application) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(a.LoadSession)

	if a.Debug {
		router.Use(middleware.Logger)
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		a.Session.Put(r.Context(), "test", "John Cena")

		err := a.Render(w, r, "index", nil)
		if err != nil {
			log.Fatalln(err)
		}
	})

	router.Get("/comments", func(w http.ResponseWriter, r *http.Request) {
		vars := make(jet.VarMap)
		vars.Set("test", a.Session.GetString(r.Context(), "test"))

		err := a.Render(w, r, "index", vars)
		if err != nil {
			log.Fatalln(err)
		}
	})

	return router
}
