package server

import (
	"github.com/justinas/nosurf"
	"net/http"
)

func (a *Application) LoadSession(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}

func (a *Application) authRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := a.Session.GetInt(r.Context(), sessionKeyUserId)
		if userId == 0 {
			http.Redirect(w,r,"login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cash-Control", "no-store")
		next.ServeHTTP(w, r)
	}
}

func (a *Application) CSRFTokenRequired(next http.Handler) http.Handler {
	handler := nosurf.New(next)

	return handler
}