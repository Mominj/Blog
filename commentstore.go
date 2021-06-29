package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID        int32  `db:"id"`
	BlogID    int    `db:"blog_id"`
	UserID    int32  `db:"user_id"`
	Comment   string `db:"comment"`
}
type CommentTemplate struct{
	Comment
	Errors map[string]error
}
func commnetStore(res http.ResponseWriter,req *http.Request){
	err := req.ParseForm()
	if err != nil {
		log.Println("error while parsing form: ", err.Error())
	}

	var comment Comment

	err = decoder.Decode(&comment, req.PostForm)
	if err != nil {
		log.Println("error while decoding form to struct: ", err.Error())
	}
	vars := mux.Vars(req)

	blogId, err := strconv.Atoi(vars["id"])
	if err!=nil{
		log.Println(err.Error())
	} 

    comment.BlogID = blogId

    session, err := SessionStore.Get(req, "first_app-session")
	if err!=nil{
		log.Println(err)
	}
	comment.UserID = session.Values["user_id"].(int32)

	
	query := `
		INSERT INTO comment(
			blog_id,
			user_id,
			comment
		)
		VALUES(
			:blog_id,
			:user_id,
			:comment
		)
		RETURNING id
	`

	var id int32
	stmt, err := DB.PrepareNamed(query)

	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, comment); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}

	//http.Redirect(res, req, "/blogs/{id}", http.StatusSeeOther)
	http.Redirect(res, req, fmt.Sprintf("/blogs/%d", blogId), http.StatusSeeOther)
}