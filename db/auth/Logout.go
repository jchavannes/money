package auth

import (
	"errors"
	"fmt"
	"git.jasonc.me/main/money/db"
)

func Logout(cookieId string) error {
	session, err := db.GetSession(cookieId)
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting session: %s\n", err))
	}

	session.HasLoggedOut = true
	err = session.Save()
	if err != nil {
		return errors.New(fmt.Sprintf("Error saving session: %s\n", err))
	}

	return nil
}
