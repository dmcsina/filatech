package data

import "database/sql"

type Models struct {
	Users UserModel
}

func newModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
