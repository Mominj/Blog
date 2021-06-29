package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type SingleBlogs struct {
	ID          int       `db:"id"`
	UserID      int       `db:"userid"`
	Title       string    `db:"title"`
	Message     string    `db:"message"`
	CreateTime  time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}
type BlogComments struct {
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Comment     string    `db:"comment"`
	CreateTime  time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}
type UserName struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type SingleBlogsTemplate struct {
	SingleBlogs
	UserName
	BlogComments
	Errors          map[string]error
	BlogCommentlist []BlogComments
}

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=password dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}

	SessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

}

func singleBlogShow(res http.ResponseWriter, req *http.Request) {
	parsedTemplate, err := template.ParseFiles("templates/user/singleblogdetails.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}
	var singleblogs SingleBlogs
	vars := mux.Vars(req)

	query := `SELECT  id, userid, title, message , created_at, updated_at FROM blogss where id=$1`

	if err := DB.Get(&singleblogs, query, vars["id"]); err != nil {
		if err != nil {
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
		return
	}
	var blogComments []BlogComments

	quer := `SELECT  userss.first_name, userss.last_name, comment.comment, comment.created_at, comment.updated_at FROM userss INNER JOIN comment on comment.user_id = userss.id where blog_id=$1`

	if err := DB.Select(&blogComments, quer, vars["id"]); err != nil {
		if err != nil {
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
		return
	}
	var userName UserName
	userName.ID = singleblogs.UserID

	que := `SELECT  first_name, last_name FROM userss where id=$1`

	if err := DB.Get(&userName, que, userName.ID); err != nil {
		if err != nil {
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, SingleBlogsTemplate{
		SingleBlogs:     singleblogs,
		BlogCommentlist: blogComments,
		UserName:        userName,
	}); err != nil {
		log.Println("error while executing template: ", err.Error())
		return
	}

}
