package main

import (
	"log"
	"net/http"
)

func logout(res http.ResponseWriter, req *http.Request) {
	session, err := SessionStore.Get(req, "first_app-session")
	if err != nil {
		log.Println("error while getting session: ", err.Error())
		return
	}

	session.Values["user_id"] = ""

	if err := session.Save(req, res); err != nil {
		log.Println("error while saving session: ", err.Error())
		return
	}

	http.Redirect(res, req, "/blogsdetails/1", http.StatusSeeOther)
}
