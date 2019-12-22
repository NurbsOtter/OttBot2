package models

import (
	"OttBot2/settings"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// MakeDB connects to a database based on the dbConnectString
func MakeDB(dbConnectString string) error {
	newDb, err := sql.Open("mysql", dbConnectString)
	if err != nil {
		return err
	}

	db = newDb
	timeout, err := time.ParseDuration(settings.GetDatabaseTimeout())
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(timeout)
	return nil
}
