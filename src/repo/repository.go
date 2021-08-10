package repo

import "database/sql"

type SqlFunctions interface {
}

type DBRepo struct {
	DB *sql.DB
}
