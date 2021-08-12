package repo

import (
	"context"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// DBRepoHMAC holding our db functionalities schema we need to implement
type DBRepoHMAC interface {
	CreateUserDataModelMethod() error
	SaveUserWithHAMCMethod(user *DataModel) error
	GetUserByItsEmailHMACMethod(um string) (string, error)
}

type UserClaims struct {
	*jwt.StandardClaims
	Email string
}

// DataModel for storing user stuff
type DataModel struct {
	ID        int    `json:"id"`
	Mail      string `json:"mail"`
	Key       string `json:"key,omitempty"`
	HmacToken string `json:"hmac_token,omitempty"`
}

// DBRepo use for holding *sql.DB
type DBRepo struct {
	DB *sql.DB
}

// CreateUserDataModelMethod use for creating DataModel for every user
func (dbr *DBRepo) CreateUserDataModelMethod() error {
	err := dbr.PingingDB()
	if err != nil {
		log.Println(err.Error() + "; error occurred during pinging db.")
		return err
	}

	statement := `CREATE TABLE users_hmac (user_id integer primary key, user_mail varchar(255), user_key varchar(511))`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*9)
	defer cancel()
	_, err = dbr.DB.ExecContext(ctx, statement)
	if err != nil {
		log.Println(err.Error() + "; error in creating table")
		return err
	}

	return nil
}

// SaveUserWithHAMCMethod use for save user with its associated key into Database
func (dbr *DBRepo) SaveUserWithHAMCMethod(user *DataModel) error {
	err := dbr.PingingDB()
	if err != nil {
		log.Println(err.Error() + "; error occurred during pinging db.")
		return err
	}

	query := `INSERT INTO users_hmac (user_id, user_mail, user_key) VALUES 
(?, ?, ?)`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*9)
	defer cancel()
	_, err = dbr.DB.ExecContext(ctx, query,
		user.ID,
		user.Mail,
		user.Key,
	)
	if err != nil {
		log.Println(err.Error() + "; error in inserting a user")
		return err
	}

	return nil
}

// GetUserByItsEmailHMACMethod use for getting user information by its email
func (dbr *DBRepo) GetUserByItsEmailHMACMethod(um string) (string, error) {
	err := dbr.PingingDB()
	if err != nil {
		log.Println(err.Error() + "; error occurred during pinging db.")
		return "", err
	}

	query := `SELECT user_key from users_hmac WHERE user_mail=?`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var userKey string = ""
	row := dbr.DB.QueryRowContext(ctx, query, um)
	err = row.Scan(&userKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return userKey, nil
}
