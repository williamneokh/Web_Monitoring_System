package handler

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"williamGoInAction1/dataStructure"
)

func GetUser(res http.ResponseWriter, req *http.Request) dataStructure.User {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}

	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var myUser dataStructure.User
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

func AlreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username]
	return ok
}
