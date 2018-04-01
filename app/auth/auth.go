package auth

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/money/app/db"
)

func IsLoggedIn(cookieId string) bool {
	session, err := db.GetSession(cookieId)
	if err != nil {
		return false
	}
	if session.UserId > 0 && ! session.HasLoggedOut {
		return true
	}
	return false
}

func GetSessionUser(cookieId string) (*db.User, error) {
	session, err := db.GetSession(cookieId)
	if err != nil || session.UserId == 0 || session.HasLoggedOut {
		return nil, jerr.Get("Error getting session", err)
	}
	user, err := db.GetUserById(session.UserId)
	if err != nil {
		return nil, jerr.Get("Error getting session user", err)
	}
	return user, nil
}
