package urlHandlers

import (
	"backend/helpers"
	"backend/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login attempt!")

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// check user auth
	success, userCookie := validators.ValidateUserLogin(email, password)
	var callback = make(map[string]string)
	if success {
		sessionCookie := http.Cookie{
			Name:     "socialNetworkSession",
			Value:    userCookie,
			Expires:  time.Now().Add(time.Minute * 10),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &sessionCookie)

		authCookie := http.Cookie{
			Name:     "socialNetworkAuth",
			Value:    "true",
			Expires:  time.Now().Add(time.Minute * 10),
			Path:     "/",
		}
		http.SetCookie(w, &authCookie)
		callback["login"] = "success"
	} else {
		callback["login"] = "fail"
		callback["error"] = userCookie
	}
	// fmt.Println(email)
	// fmt.Println(password)

	writeData, err := json.Marshal(callback)
	helpers.CheckErr("handleLogin", err)
	w.Write(writeData)
}
