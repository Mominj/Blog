package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)
type SinBlog struct {
	ID int `db:"id"`
	UserID int `db:"userid"`
	Title string `db:"title"`
	Message string `db:"message"`
	Category_id int `db:"category_id"`
	CreateTime time.Time `db:"created_at"`
    UpdatedTime time.Time `db:"updated_at"`
}
type StoreSingleBlog struct{
	ID int `db:"id"`
	Message string `db:"message"`
}
type SinBlogTemplate struct{
	SinBlog
	SinBloglist SinBlog
	Errors   map[string]error
}
func blogEdit(res http.ResponseWriter,req *http.Request) {
	parsedTemplate, err := template.ParseFiles(
		"templates/user/singleblogedit.html",
	)
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
	return
   }
	vars := mux.Vars(req)

	blogid, err :=  strconv.Atoi(vars["id"])
	if err!=nil{
		fmt.Println("error occurs when convert string to int",err.Error())
	}
	/*session, err := SessionStore.Get(req, "first_app-session")
			if err != nil {
				log.Println("error while getting session: ", err.Error())
			}
	userID := session.Values["user_id"]*/
	
	var sinBlog SinBlog
	query := `
		SELECT * FROM blogss
		WHERE id = $1
	`
	if err := DB.Get(&sinBlog, query, blogid); err != nil {
		log.Println("error while getting blog from db: ", err.Error())
	}
	if err := parsedTemplate.Execute(res, SinBlogTemplate{
		SinBloglist: sinBlog,	
	}); err != nil {
		log.Println("error while executing template: ", err.Error())
	 }
	
}



func blogUpdate(res http.ResponseWriter,req *http.Request){
	err := req.ParseForm()
	if err != nil {
	log.Println("error while parsing form: ", err.Error())
}

var storeSingleBlog StoreSingleBlog

	err = decoder.Decode(&storeSingleBlog,req.PostForm)
if err != nil {
	log.Println("error while decoding form to struct: ", err.Error())
}
	vars := mux.Vars(req)

	blogid, err :=  strconv.Atoi(vars["id"])
	if err!=nil{
		fmt.Println("error occurs when convert string to int",err.Error())
}
 storeSingleBlog.ID  = blogid
//log.Println("storeSingleBlog.ID is = ",storeSingleBlog.ID)
//log.Println("Message is = ",storeSingleBlog.Message)
query := `
UPDATE blogss
		SET  message = :message
		WHERE blogss.id = :id
`

stmt, err := DB.PrepareNamed(query)

if err != nil {
	log.Println("db error: failed prepare ", err.Error())
	return
}

if _, err := stmt.Exec(&storeSingleBlog); err != nil {
	log.Println("db error: failed to update data ", err.Error())
	return
}
	http.Redirect(res, req, fmt.Sprintf("/blogs/%d", blogid), http.StatusSeeOther)
}

func blogDelete(res http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	blogId := vars["id"]
	var sinBlog SinBlog
	blogID, err := strconv.Atoi(blogId)
	if err != nil {
		log.Println("error while convert string to int : ", err.Error())
	}
	sinBlog.ID = blogID
	log.Println("blog id",blogID)
	query := `
		DELETE 
		FROM blogss
		WHERE id = :id
	`
	stmt, err := DB.PrepareNamed(query)
	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if _, err := stmt.Exec(&sinBlog); err != nil {
		log.Println("db error: failed to update data ", err.Error())
		return
	}

	http.Redirect(res, req, "/blogsdetails/1", http.StatusSeeOther)
}