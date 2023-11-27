package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Application) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	if a.Debug {
		router.Use(middleware.Logger)
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := a.Render(w, r, "index", nil); err != nil {
			log.Fatalln(err)
		}
	})

	return router
}
