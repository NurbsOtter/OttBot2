package models

import (
	"database/sql"
	"time"
)

type Warning struct {
	ID          int64
	WarningText string
	WarnDate    time.Time
	TgID        int64
}

// AddWarningToId logs a new warning for a given user
func AddWarningToID(tgID int64, warnText string) error {
	_, err := db.Exec("INSERT INTO warning(warningText,warnDate,tgID) VALUES(?,?,?)", warnText, time.Now(), tgID)
	return err
}

// GetUsersWarnings finds all warnings for a given user
func GetUsersWarnings(userIn *ChatUser) ([]Warning, error) {
	stmt, err := db.Prepare("SELECT warningText,warnDate FROM warning WHERE tgID = ?")
	if err != nil {
		return nil, err
	}

	var outWarns []Warning
	rows, err := stmt.Query(userIn.TgID)
	// User not found
	switch {
	case err == sql.ErrNoRows:
		return outWarns, nil

	// Unknown error
	case err != nil:
		return nil, err

	// User found
	default:
		for rows.Next() {
			newWarn := Warning{}
			rows.Scan(&newWarn.WarningText, &newWarn.WarnDate)
			outWarns = append(outWarns, newWarn)
		}
		return outWarns, nil
	}
}
