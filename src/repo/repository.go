package repo

import "database/sql"

type DBRepoFunctions interface {
}

type DataModel struct {
	ID        int    `json:"id"`
	Mail      string `json:"mail"`
	Key       string `json:"key,omitempty"`
	HmacToken string `json:"hmac_token,omitempty"`
}

type DBRepo struct {
	DB *sql.DB
}
