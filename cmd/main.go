package main

import (
	"log"
	"os"

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

	if err := app.StartServer(); err != nil {
		log.Fatal(err)
	}
}
