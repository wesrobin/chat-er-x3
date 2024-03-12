package main

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "schema.sql")
	if err != nil {
		return nil, err
	}

	return db, nil
}
