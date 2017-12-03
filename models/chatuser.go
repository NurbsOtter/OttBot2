package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type ChatUser struct {
	ID               int64
	UserName         string
	TgID             int64
	PingAllowed      bool
	MarkovAskAllowed bool
}
type UserAlias struct {
	ID         int64
	Name       string
	UserID     int64
	ChangeDate time.Time
}

func ChatUserFromID(inID int64) *ChatUser {
	newUser := &ChatUser{}
	err := db.QueryRow("SELECT id,userName,tgID,pingAllowed,MarkovAskAllowed FROM chatUser WHERE id = ?", inID).Scan(&newUser.ID, &newUser.UserName, &newUser.TgID, &newUser.PingAllowed, &newUser.MarkovAskAllowed)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return newUser
	}
}
func ChatUserFromTGID(tgID int, userName string) *ChatUser {
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed,MarkovAskAllowed FROM chatUser WHERE tgID = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundUser := &ChatUser{}
	err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed, &foundUser.MarkovAskAllowed)
	switch {
	case err == sql.ErrNoRows:
		insStmt, err := db.Prepare("INSERT INTO chatUser(userName,tgID) VALUES(?,?)")
		if err != nil {
			panic(err)
		}
		defer stmt.Close()
		insStmt.Exec(userName, tgID)
		err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed, &foundUser.MarkovAskAllowed)
		return foundUser
	case err != nil:
		panic(err)
	default:
	}
	if foundUser.UserName != userName { //They've changed their username! Update it!
		_, err := db.Exec("UPDATE chatUser SET userName = ? WHERE tgID = ?", userName, foundUser.TgID)
		if err != nil {
			fmt.Println("Failed to change tgID " + string(foundUser.TgID))
		}
		foundUser.UserName = userName
	}
	return foundUser
}
func ChatUserFromTGIDNoUpd(tgID int) *ChatUser {
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed,MarkovAskAllowed FROM chatUser WHERE tgID = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundUser := &ChatUser{}
	err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed, &foundUser.MarkovAskAllowed)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
	}
	return foundUser
}
func UpdateAliases(firstName string, lastName string, userID int64) {
	insName := firstName + " " + lastName
	stmt, err := db.Prepare("SELECT name FROM aliases WHERE name = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundAlias := &UserAlias{}
	err = stmt.QueryRow(insName).Scan(&foundAlias.Name)
	switch {
	case err == sql.ErrNoRows:
		insStmt, err := db.Prepare("INSERT INTO aliases(name,userID,changeDate) VALUES (?,?,?)")
		if err != nil {
			panic(err)
		}
		defer insStmt.Close()
		insStmt.Exec(insName, userID, time.Now())
	case err != nil:
		panic(err)
	default:
		updateStmt, err := db.Prepare("UPDATE aliases SET changeDate = ? WHERE name = ?")
		if err != nil {
			panic(err)
		}
		defer updateStmt.Close()
		updateStmt.Exec(time.Now(), insName)
	}
}

func SearchUserByUsername(userName string) *ChatUser {
	fmt.Println(userName)
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed,MarkovAskAllowed FROM chatUser WHERE lower(userName) = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundUser := &ChatUser{}
	err = stmt.QueryRow(userName).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed, &foundUser.MarkovAskAllowed)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
	}
	fmt.Println(foundUser)
	return foundUser
}
func GetLatestAliasFromUserID(userID int64) *UserAlias {
	stmt, err := db.Prepare("SELECT id,name,userID,changeDate FROM aliases WHERE userID = ? ORDER BY changeDate DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundAlias := &UserAlias{}
	err = stmt.QueryRow(userID).Scan(&foundAlias.ID, &foundAlias.Name, &foundAlias.UserID, &foundAlias.ChangeDate)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
	}
	return foundAlias

}
func GetAliases(u *ChatUser) []UserAlias {
	rows, err := db.Query("SELECT id,name,userID,changeDate FROM aliases WHERE userID = ?", u.ID)
	var outAliases []UserAlias
	switch {
	case err == sql.ErrNoRows:
		return outAliases
	case err != nil:
		panic(err)
	default:
	}
	for rows.Next() {
		newAlias := UserAlias{}
		rows.Scan(&newAlias.ID, &newAlias.Name, &newAlias.UserID, &newAlias.ChangeDate)
		outAliases = append(outAliases, newAlias)
	}
	return outAliases
}
func LookupAlias(query string) []ChatUser {
	newQuery := "%" + strings.ToLower(query) + "%"
	rows, err := db.Query("SELECT DISTINCT chatUser.id,chatUser.userName,chatUser.tgID,chatUser.pingAllowed,chatUser.MarkovAskAllowed FROM chatUser JOIN aliases ON aliases.userID = chatUser.id WHERE lower(aliases.name) LIKE ? LIMIT 20", newQuery)
	var outUsers []ChatUser
	switch {
	case err == sql.ErrNoRows:
		return outUsers
	case err != nil:
		panic(err)
	default:
	}
	defer rows.Close()
	for rows.Next() {
		newUser := ChatUser{}
		rows.Scan(&newUser.ID, &newUser.UserName, &newUser.TgID, &newUser.PingAllowed, &newUser.MarkovAskAllowed)
		if newUser.UserName == "" {
			newUser.UserName = "None"
		}
		outUsers = append(outUsers, newUser)
	}
	return outUsers
}
func SearchAliases(query string) []UserAlias {
	newQuery := "%" + strings.ToLower(query) + "%"
	rows, err := db.Query("SELECT id,name,userID,changeDate FROM aliases WHERE name LIKE ?", newQuery)
	var outAliases []UserAlias
	switch {
	case err == sql.ErrNoRows:
		return outAliases
	case err != nil:
		panic(err)
	default:
	}
	for rows.Next() {
		newAlias := UserAlias{}
		rows.Scan(&newAlias.ID, &newAlias.Name, &newAlias.UserID, &newAlias.ChangeDate)
		outAliases = append(outAliases, newAlias)
	}
	return outAliases
}

//func (u *ChatUser) ToggleModPing() { //Refactor this for use via command
//	newRights := !u.PingAllowed
//	_, err := db.Exec("UPDATE chatUser SET pingAllowed=? WHERE id=?", newRights, u.ID)
//	if err != nil {
//		panic(err)
//	}
//}
func SetModPing(userId int64, status bool) {
	_, err := db.Exec("UPDATE chatUser SET pingAllowed=? WHERE id=?", status, userId)
	if err != nil {
		panic(err)
	}
}
func SetMarkovUse(userId int64, status bool) {
	_, err := db.Exec("UPDATE chatUser SET MarkovAskAllowed=? WHERE id=?", status, userId)
	if err != nil {
		panic(err)
	}
}
