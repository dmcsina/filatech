package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
}

func main() {

	app := application{}
	app.setEnvironmentVariables("hide.env")

	var user string
	var database string
	var password string
	var sslmode string

	flag.StringVar(&user, "db-user", os.Getenv("user"), "DB user")
	flag.StringVar(&database, "db", os.Getenv("dbname"), "DB")
	flag.StringVar(&password, "password", os.Getenv("password"), "DB user password")
	flag.StringVar(&sslmode, "sslmode", os.Getenv("sslmode"), "Ssl mode")

	var dsn string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, database, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully established connection with the database")
	defer db.Close()

	server := http.Server{
		Addr:         ":8080",
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
