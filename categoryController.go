package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type Category struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	CreateTime  time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}
type CategoryTemplate struct {
	Category
	Errors map[string]error
}

func blogsCategoryCreate(res http.ResponseWriter, req *http.Request) {
	parsedTemplate, err := template.ParseFiles("templates/user/categorycreate.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}
	var categoryTemplate CategoryTemplate
	if err := parsedTemplate.Execute(res, categoryTemplate); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
}

func blogsCategoryStore(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error while parsing form: ", err.Error())
	}

	var categoryTemplate CategoryTemplate

	err = decoder.Decode(&categoryTemplate, req.PostForm)
	if err != nil {
		log.Println("error while decoding form to struct: ", err.Error())
	}

	query := `
	INSERT INTO categories(
		name
		)
		VALUES(
		:name
		)
		RETURNING id
    `

	var id int32
	stmt, err := DB.PrepareNamed(query)

	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, categoryTemplate); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}
	http.Redirect(res, req, "/blogsdetails/1", http.StatusSeeOther)
	//http.Redirect(res, req,"/blogs", http.StatusSeeOther)
	//http.Redirect(res, req, fmt.Sprintf("blogsdetails/%d", 1), http.StatusSeeOther)
}



