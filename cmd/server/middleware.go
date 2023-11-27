package server

import "net/http"

func (a *Application) LoadSession(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}