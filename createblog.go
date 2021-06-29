package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	ID          int32  `db:"id"`
	UserID      int32  `db:"userid"`
	Title       string `db:"title"`
	Message     string `db:"message"`
	Category_id int32  `db:"category_id"`
}

type UserStoreTemplate struct {
	UserStore
	Errors map[string]error
}

func createblog(res http.ResponseWriter, req *http.Request) {
	ViewUserStoreTemplate(res, UserStoreTemplate{})
}

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=password dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}

	SessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

}

func blogStore(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error while parsing form: ", err.Error())
	}

	var userstore UserStore

	err = decoder.Decode(&userstore, req.PostForm)
	if err != nil {
		log.Println("error while decoding form to struct: ", err.Error())
	}
	session, err := SessionStore.Get(req, "first_app-session")

	if err != nil {
		log.Println("error while session store : ", err.Error())
	}
	userstore.UserID = session.Values["user_id"].(int32)

	query := `
	INSERT INTO blogss(
	userid,
	title,
	message,
	category_id
	)
	VALUES(
	:userid,
	:title,
	:message,
	:category_id
	)
	RETURNING id
	`

	var id int32
	stmt, err := DB.PrepareNamed(query)

	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, userstore); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}
	//http.Redirect(res, req, "/home", http.StatusSeeOther)
	//http.Redirect(res, req,"/blogs", http.StatusSeeOther)
	http.Redirect(res, req, fmt.Sprintf("blogsdetails/%d", 1), http.StatusSeeOther)
}

func ViewUserStoreTemplate(res http.ResponseWriter, userstoreTemplate UserStoreTemplate) {
	parsedTemplate, err := template.ParseFiles("templates/user/createblog.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, userstoreTemplate); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
}
