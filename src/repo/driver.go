package repo

import (
	"context"
	"database/sql"
	"log"
)

var Repo *DBRepo

func NewDB(dsn string) (*DBRepo, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	Repo = &DBRepo{
		DB: db,
	}
	err = Repo.PingingDB()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return Repo, nil
}

func (dbr *DBRepo) PingingDB() error {
	err := dbr.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (dbr *DBRepo) CreateStatement(stat string, ctx context.Context) (*sql.Stmt, error) {
	statement, err := dbr.DB.PrepareContext(ctx, stat)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return statement, err
}
