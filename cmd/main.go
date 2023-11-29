package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/codepnw/web-news/cmd/server"
	"github.com/codepnw/web-news/database"
	"github.com/codepnw/web-news/models"
	"github.com/joho/godotenv"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"

	_ "github.com/lib/pq"
)

func main() {
	migrate := flag.Bool("migrate", false, "should migrate - drop all tables")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading .env file failed.")
	}

	dbs, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbs.Close()

	upper, err := postgresql.New(dbs)
	if err != nil {
		log.Fatal(err)
	}
	defer func(upper db.Session) {
		err := upper.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(upper)

	if *migrate {
		fmt.Println("running migrations...")
		err = database.MigrateDB(upper)
		if err != nil {
			log.Fatal(err)
		}
	}

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
		Model:   models.New(upper),
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
	app.Session.Store = postgresstore.New(dbs)

	if err := app.StartServer(); err != nil {
		log.Fatal(err)
	}
}
