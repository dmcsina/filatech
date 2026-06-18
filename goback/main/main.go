package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"najaftech.filatech/data"
)

type application struct {
	logger    *log.Logger
	UserModel data.UserModel
	JwtSecret []byte
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := application{
		logger: logger,
	}
	app.setEnvironmentVariables("hide.env")

	var user string
	var database string
	var password string
	var sslmode string
	var secretkey string
	var jwtsecretString string

	flag.StringVar(&user, "db-user", os.Getenv("user"), "DB user")
	flag.StringVar(&database, "db", os.Getenv("dbname"), "DB")
	flag.StringVar(&password, "password", os.Getenv("password"), "DB user password")
	flag.StringVar(&sslmode, "sslmode", os.Getenv("sslmode"), "Ssl mode")
	flag.StringVar(&secretkey, "secretkey", os.Getenv("secretkey"), "secretKey")
	flag.StringVar(&jwtsecretString, "jwtsecretkey", os.Getenv("jwtsecret"), "jwtsecret")
	var dsn string = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, database, sslmode)

	app.JwtSecret = []byte(jwtsecretString)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully established connection with the database")
	defer db.Close()

	app.UserModel = *data.NewUserModel(db, secretkey)

	server := http.Server{
		Addr:         ":8080",
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
