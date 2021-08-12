package repo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Repo *DBRepo

// NewDB create new sqlite3 db for use by its name
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

// PingingDB use for ping db
func (dbr *DBRepo) PingingDB() error {
	err := dbr.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

// CreateStatement use for creating *sql.Stmt for a specific query or statement
func (dbr *DBRepo) CreateStatement(stat string, ctx context.Context) (*sql.Stmt, error) {
	statement, err := dbr.DB.PrepareContext(ctx, stat)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return statement, err
}

// DisposeFunction use for dispose sql.DB
func (dbr *DBRepo) DisposeFunction() error {
	err := dbr.DB.Close()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}
