package validators

import (
	"backend/database"
	"backend/helpers"

	"github.com/google/uuid"
)

// check if username or email exists in users table
func ValidateUserRegistration(email, username string) (bool, bool) {
	email = helpers.StandardizeName(email)
	emailExists := database.GetEmailIfExists(email)
	if len(username) > 0 {
		username = helpers.StandardizeName(username)
		usernameExists := database.GetUsernameIfExists(username)
		return emailExists, usernameExists
	}
	return emailExists, false
}

func ValidateUserLogin(email, password string) (bool, string) {
	email = helpers.StandardizeName(email)
	userId, userPsw := database.GetUserIdPswByEmail(email)
	if userId != "0" {
		if helpers.CheckPassword(password, userPsw) {
			hash := uuid.New().String()
			// set user and UUID in DB
			database.SetToSessions(userId, hash)
			return true, hash
		} else {
			return false, "Password error"
		}
	} else {
		return false, "Email error"
	}
}

func ValidateUserSession(cookie string) bool {
	return database.GetUserSession(cookie)
}
