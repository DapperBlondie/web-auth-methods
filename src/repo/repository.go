package repo

import "database/sql"

type DBRepo struct {
	DB *sql.DB
}
