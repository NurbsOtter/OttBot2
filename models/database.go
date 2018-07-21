package models

import (
	"OttBot2/settings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func MakeDB(dbConnectString string) {
	newDb, err := sql.Open("mysql", dbConnectString)
	if err != nil {
		panic(err)
	}
	db = newDb
	timeout, err := time.ParseDuration(settings.GetDatabaseTimeout())
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(timeout)
}
