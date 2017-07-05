package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func MakeDB(dbConnectString string) {
	newDb, err := sql.Open("mysql", dbConnectString)
	if err != nil {
		panic(err)
	}
	db = newDb
}
