package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Application) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(a.CSRFTokenRequired)
	router.Use(a.LoadSession)

	if a.Debug {
		router.Use(middleware.Logger)
	}

	router.Get("/", a.homeHandler)
	router.Get("/comments/{postId}", a.commentHandler)
	router.Post("/comments/{postId}", a.commentPostHandler)

	router.Get("/login", a.loginHandler)
	router.Post("/login", a.loginPostHandler)
	router.Get("/signup", a.signUpHandler)
	router.Post("/signup", a.signUpPostHandler)
	router.Get("/logout", a.authRequired(a.logoutHandler))

	router.Get("/vote", a.authRequired(a.voteHandler))
	router.Get("/submit", a.authRequired(a.submitHandler))
	router.Post("/submit", a.authRequired(a.submitPostHandler))


	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return router
}
