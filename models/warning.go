package models

import (
	"database/sql"
	"time"
)

type Warning struct {
	ID          int64
	UserID      int64
	WarningText string
	WarnDate    time.Time
}

func AddWarningToID(userID int64, warnText string) {
	db.Exec("INSERT INTO warning(userID,warningText,warnDate) VALUES(?,?,?)", userID, warnText, time.Now())
}

func GetUsersWarnings(userIn *ChatUser) []Warning {
	stmt, err := db.Prepare("SELECT warningText,warnDate FROM warning WHERE userID = ?")
	if err != nil {
		panic(err)
	}
	var outWarns []Warning
	rows, err := stmt.Query(userIn.ID)
	if err == sql.ErrNoRows {
		return outWarns
	} else {
		for rows.Next() {
			newWarn := Warning{}
			rows.Scan(&newWarn.WarningText, &newWarn.WarnDate)
			outWarns = append(outWarns, newWarn)
		}
		return outWarns
	}
}
