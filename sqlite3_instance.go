package main

import (
	"database/sql"
)

var db *sql.DB

func openDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./remote_config.db")
	if err != nil {
		return err
	}
	return nil
}
