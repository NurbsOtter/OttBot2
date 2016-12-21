package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func MakeDB(dbFile string) {
	newDb, err := sql.Open("sqlite3", "file:"+dbFile+"?loc=auto")
	if err != nil {
		panic(err)
	}
	db = newDb
}
