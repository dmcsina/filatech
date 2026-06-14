package main

import (
	"net/http"
	"time"
)

type application struct {
}

func main() {

	app := application{}

	server := http.Server{
		Addr:         ":8080",
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
