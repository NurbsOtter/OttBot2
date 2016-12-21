package models

import (
	"database/sql"
	"github.com/kataras/iris"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	UserName string
	Password string
}

func InsertUser(userName string, password string) bool {
	stmt, err := db.Prepare("INSERT INTO adminUsers(userName,password) VALUES (?,?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	if !IsUserTaken(userName) {
		newPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 11)
		stmt.Exec(userName, string(newPassword))
		return true
	} else {
		return false
	}
}

//makes a user from the session variables.
func GetUserFromContext(ctx *iris.Context) *User {
	userID, err := ctx.Session().GetInt("userID")
	if err != nil {
		return nil
	}
	if userID < 0 {
		return nil
	}
	stmt, err := db.Prepare("SELECT * FROM adminUsers WHERE id = ?")
	if err != nil {
		panic(err)
	}
	foundUser := &User{}
	err = stmt.QueryRow(userID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
	}
	return foundUser
}
func VerifyUser(userName string, password string) *User {
	stmt, err := db.Prepare("SELECT * FROM adminUsers WHERE userName = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	foundUser := &User{}
	err = stmt.QueryRow(userName).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.Password)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	switch {
	case err == bcrypt.ErrMismatchedHashAndPassword:
		return nil
	case err != nil:
		panic(err)
	default:
		return foundUser
	}
}
func IsUserTaken(userName string) bool {
	stmt, err := db.Prepare("SELECT userName FROM adminUsers WHERE userName = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	var foundUser string
	err = stmt.QueryRow(userName).Scan(&foundUser)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		panic(err)
	default:
		return true
	}
}
