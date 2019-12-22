package models

import (
	"database/sql"
	"strings"
	"time"
)

type ChatUser struct {
	ID          int64
	UserName    string
	TgID        int64
	PingAllowed bool
}
type UserAlias struct {
	ID         int64
	Name       string
	UserID     int64
	ChangeDate time.Time
}

// ChatUserFromID finds a user with the provided ID
func ChatUserFromID(inID int64) (*ChatUser, error) {
	newUser := &ChatUser{}

	err := db.QueryRow("SELECT id,userName,tgID,pingAllowed FROM chatUser WHERE id = ?", inID).Scan(&newUser.ID, &newUser.UserName, &newUser.TgID, &newUser.PingAllowed)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// ChatUserFromTGID finds a user with the provided TGID
func ChatUserFromTGID(tgID int, userName string) (*ChatUser, error) {
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed FROM chatUser WHERE tgID = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	foundUser := &ChatUser{}
	err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed)

	switch {
	// User not found
	case err == sql.ErrNoRows:
		insStmt, err := db.Prepare("INSERT INTO chatUser(userName,tgID) VALUES(?,?)")
		if err != nil {
			return nil, err
		}

		defer stmt.Close()
		insStmt.Exec(userName, tgID)
		err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed)
		if err != nil {
			return nil, err
		}

		return foundUser, nil

	// Unknown error
	case err != nil:
		return nil, err

	// User found
	default:
		// Username update detection
		if foundUser.UserName != userName {
			_, err := db.Exec("UPDATE chatUser SET userName = ? WHERE tgID = ?", userName, foundUser.TgID)
			if err != nil {
				return nil, err
			}

			foundUser.UserName = userName
		}
		return foundUser, nil
	}
}

// ChatUserFromTGIDNoUpd finds a user with the provided TGID without updating the database
func ChatUserFromTGIDNoUpd(tgID int) (*ChatUser, error) {
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed FROM chatUser WHERE tgID = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	foundUser := &ChatUser{}

	err = stmt.QueryRow(tgID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

// UpdateAliases records new names users adopt as aliases, and updates the time changed if a user adopts an old alias
func UpdateAliases(firstName string, lastName string, userID int64) error {
	insName := firstName + " " + lastName
	stmt, err := db.Prepare("SELECT name FROM aliases WHERE name = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	foundAlias := &UserAlias{}

	err = stmt.QueryRow(insName).Scan(&foundAlias.Name)
	switch {
	// User not Found
	// Add the current name as an alias, and record the id and time
	case err == sql.ErrNoRows:
		insStmt, err := db.Prepare("INSERT INTO aliases(name,userID,changeDate) VALUES (?,?,?)")
		if err != nil {
			return err
		}

		defer insStmt.Close()
		insStmt.Exec(insName, userID, time.Now())

		return nil

	// Unknown error
	case err != nil:
		return err

	// User found
	// The user has returned to a previous alias; update the time changed
	default:
		updateStmt, err := db.Prepare("UPDATE aliases SET changeDate = ? WHERE name = ?")
		if err != nil {
			return err
		}

		defer updateStmt.Close()
		updateStmt.Exec(time.Now(), insName)
		return nil
	}
}

// SearchUserByUsername finds a user with the provided username
func SearchUserByUsername(userName string) (*ChatUser, error) {
	stmt, err := db.Prepare("SELECT id,userName,tgID,pingAllowed FROM chatUser WHERE lower(userName) = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	foundUser := &ChatUser{}

	err = stmt.QueryRow(userName).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.TgID, &foundUser.PingAllowed)
	switch {
	// User not found
	case err == sql.ErrNoRows:
		return nil, err

	// Unknown error
	case err != nil:
		return nil, err

	// User found
	default:
		return foundUser, nil
	}
}

// GetLatestAliasFromUserID finds the most recently used alias for a given userID
func GetLatestAliasFromUserID(userID int64) (*UserAlias, error) {
	stmt, err := db.Prepare("SELECT id,name,userID,changeDate FROM aliases WHERE userID = ? ORDER BY changeDate DESC LIMIT 1")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	foundAlias := &UserAlias{}

	err = stmt.QueryRow(userID).Scan(&foundAlias.ID, &foundAlias.Name, &foundAlias.UserID, &foundAlias.ChangeDate)
	switch {
	// User not found
	case err == sql.ErrNoRows:
		return nil, err

	// Unknown error
	case err != nil:
		return nil, err

	// User found
	default:
		return foundAlias, nil
	}
}

// GetAliases finds aliases for a given ChatUser
func GetAliases(u *ChatUser) ([]UserAlias, error) {
	rows, err := db.Query("SELECT id,name,userID,changeDate FROM aliases WHERE userID = ?", u.ID)
	var outAliases []UserAlias
	switch {
	// User not found
	case err == sql.ErrNoRows:
		return outAliases, nil

	// Uknown error
	case err != nil:
		return nil, err

	// User found
	default:
		for rows.Next() {
			newAlias := UserAlias{}
			rows.Scan(&newAlias.ID, &newAlias.Name, &newAlias.UserID, &newAlias.ChangeDate)
			outAliases = append(outAliases, newAlias)
		}
		return outAliases, nil
	}
}

// LookupAlias finds a ChatUser even if only part of one of their aliases is provided
func LookupAlias(query string) ([]ChatUser, error) {
	newQuery := "%" + strings.ToLower(query) + "%"
	rows, err := db.Query("SELECT DISTINCT chatUser.id,chatUser.userName,chatUser.tgID,chatUser.pingAllowed FROM chatUser JOIN aliases ON aliases.userID = chatUser.id WHERE lower(aliases.name) LIKE ? LIMIT 20", newQuery)
	var outUsers []ChatUser
	switch {
	// User not found
	case err == sql.ErrNoRows:
		return outUsers, nil

	// Unknown error
	case err != nil:
		return nil, err

	// User found
	default:
		defer rows.Close()
		for rows.Next() {
			newUser := ChatUser{}
			rows.Scan(&newUser.ID, &newUser.UserName, &newUser.TgID, &newUser.PingAllowed)
			if newUser.UserName == "" {
				newUser.UserName = "None"
			}
			outUsers = append(outUsers, newUser)
		}
		return outUsers, nil
	}
}

// SetModPing sets the user's ability to use the moderator ping functionality
func SetModPing(userId int64, status bool) error {
	_, err := db.Exec("UPDATE chatUser SET pingAllowed=? WHERE id=?", status, userId)
	if err != nil {
		return err
	}
	return nil
}
