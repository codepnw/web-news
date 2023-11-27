package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	AppName string
	Server  Server
	Debug   bool
	ErrLog  *log.Logger
	InfoLog *log.Logger
	View    *jet.Set
	Session *scs.SessionManager
}

type Server struct {
	Host string
	Port string
	Url  string
}

func (a *Application) StartServer() error {
	host := fmt.Sprintf("%s:%s", a.Server.Host, a.Server.Port)

	srv := http.Server{
		Handler:     a.Routes(),
		Addr:        host,
		ReadTimeout: 300 * time.Second,
	}

	a.InfoLog.Printf("Server listening on :%s\n", host)

	return srv.ListenAndServe()
}
