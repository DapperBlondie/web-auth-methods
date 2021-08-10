package repo

import "database/sql"

// DBRepoFunctions holding our db functionalities schema we need to implement
type DBRepoFunctions interface {
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
