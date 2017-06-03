package auth

import "git.jasonc.me/main/money/db"

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
		return nil, err
	}
	user, err := db.GetUserById(session.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
