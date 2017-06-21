package auth

import (
	"git.jasonc.me/main/money/db"
	"golang.org/x/crypto/bcrypt"
	"github.com/jchavannes/jgo/jerr"
)

func Signup(cookieId string, username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := db.CreateUser(username, string(hashedPassword))
	if err != nil {
		return jerr.Get("Error signing up", err)
	}
	session, err := db.GetSession(cookieId)
	if err != nil {
		return jerr.Get("Error getting session", err)
	}
	session.UserId = user.Id
	err = session.Save()
	if err != nil {
		return jerr.Get("Error saving session", err)
	}
	return nil
}
