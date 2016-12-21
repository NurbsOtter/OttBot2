package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Warning struct {
	ID          int64
	UserID      int64
	WarningText string
	WarnDate    time.Time
}

func AddWarningToUsername(userName string, warnText string) {
	user := SearchUserByUsername(userName)
	fmt.Println(userName)
	if user == nil {
		fmt.Println("No luck")
		return
	}
	stmt, err := db.Prepare("INSERT INTO warning(userID,warningText,warnDate) VALUES(?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.Exec(user.ID, warnText, time.Now())
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
