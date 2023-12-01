package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
)

func (a *Application) readInt(r *http.Request, key string) int {
	v, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil {
		return 0
	}
	return v
}

func (a *Application) readIntDefault(r *http.Request, key string, dvalue int) int {
	v := a.readInt(r, key)
	if v <= 0 {
		return dvalue
	}
	return v
}

func (a *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.ErrLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}