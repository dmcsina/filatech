package data

import "database/sql"

type user struct {
	username string `json:"username"`
	email    string `json:"email"`
	password string `json:"password"`
}

type userModel struct {
	DB *sql.DB
}
