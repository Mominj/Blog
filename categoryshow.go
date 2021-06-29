package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)


type Categoryshow struct {
	ID int `db:"id"`
	UserID int `db:"userid"`
	Title string `db:"title"`
	Message string `db:"message"`
	CreateTime time.Time `db:"created_at"`
    UpdatedTime time.Time `db:"updated_at"`
	CategoryID int `db:"category_id"`
}
type CartegoryData struct{
	CID int `db:"id"`
	CName string `db:"name"`
	UpdatedTime time.Time `db:"updated_at"`
	CreateTime time.Time `db:"created_at"`
}

type CategoryshowTemplate struct{
	Categoryshowlist []Categoryshow
	CartegoryDatalist []CartegoryData
	Errors   map[string]error
}

func blogsCategoryShow(res http.ResponseWriter, req *http.Request){

		parsedTemplate, err := template.ParseFiles("templates/user/categoryshow.html")
		if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	   }

    var cartegoryData []CartegoryData
	query  := `SELECT * FROM categories`
	if err := DB.Select(&cartegoryData, query); err != nil {
		if err!=nil{
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
			return
	}
	if err := parsedTemplate.Execute(res, CategoryshowTemplate{
		CartegoryDatalist: cartegoryData,
	}); err != nil {
	log.Println("error while executing template: ", err.Error())
	return
	}	
}

func blogsCategoryShow1(res http.ResponseWriter, req *http.Request){
    
	parsedTemplate, err := template.ParseFiles("templates/user/categoryshow.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	var categoryshow []Categoryshow

	vars := mux.Vars(req)
	

   query  := `SELECT * FROM blogss where category_id=$1`

	if err := DB.Select(&categoryshow, query,vars["category_id"]); err != nil {
	if err!=nil{
		panic(err)
	}
	log.Println("error while getting user from db: ", err.Error())
		return
	}
    var cartegoryData []CartegoryData
	quer  := `SELECT * FROM categories`
	if err := DB.Select(&cartegoryData, quer); err != nil {
		if err!=nil{
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
			return
	}
	if err := parsedTemplate.Execute(res, CategoryshowTemplate{
		Categoryshowlist: categoryshow,
		CartegoryDatalist: cartegoryData,
	}); err != nil {
	log.Println("error while executing template: ", err.Error())
	return
  }
}
