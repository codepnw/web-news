package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/codepnw/web-news/cmd/server"
)

func main() {
	srv := server.Server{
		Host: "127.0.0.1",
		Port: "8000",
		Url:  "http://127.0.0.1:8000",
	}

	app := &server.Application{
		Server:  srv,
		AppName: "WebNews",
		Debug:   true,
		InfoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		ErrLog:  log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
	}

	if app.Debug {
		app.View = jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode())
	} else {
		app.View = jet.NewSet(jet.NewOSFileSystemLoader("./views"))
	}

	app.Session = scs.New()
	app.Session.Lifetime = 24 * time.Hour
	app.Session.Cookie.Persist = true
	app.Session.Cookie.Name = app.AppName
	app.Session.Cookie.Domain = app.Server.Host
	app.Session.Cookie.SameSite = http.SameSiteStrictMode

	if err := app.StartServer(); err != nil {
		log.Fatal(err)
	}
}
