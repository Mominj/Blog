package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type Blogs struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Message     string    `db:"message"`
	CreateTime  time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}

type Category1 struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	CreateTime  time.Time `db:"created_at"`
	UpdatedTime time.Time `db:"updated_at"`
}

type Paginatio struct {
	OFFSE int `db:"offset"`
	Limi  int `db:"limi"`
	Total int `db:"total"`
	Next  int `db:"next"`
	Pre   int `db:"pre"`
	Last  int `db:"last"`
}
type BlogsTemplate struct {
	Blogs
	Paginatio
	Bloglist          []Blogs
	CartegoryDatalist []Category1
	Errors            map[string]error
}

func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=password dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}

	SessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

}

func getblogs(res http.ResponseWriter, req *http.Request) {
	funcMap := template.FuncMap{
		//The name "title" is what the function will be called in the template text.
		"authCheck": func() bool {
			session, err := SessionStore.Get(req, "first_app-session")
			if err != nil {
				log.Println("error while getting session: ", err.Error())
				return false
			}
			userID := session.Values["user_id"]

			query := `SELECT EXISTS(SELECT COUNT(*) FROM userss WHERE id = $1)`
			var isExists bool
			if err := DB.Get(&isExists, query, userID); err != nil {
				log.Println("error while getting user from db: ", err.Error())
				return false
			}

			return isExists
		},
	}
	parsedTemplate, err := template.New("showblogs.html").Funcs(funcMap).ParseFiles("templates/user/showblogs.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	var paginatio Paginatio
	DB.QueryRow("SELECT count(id)  FROM blogss;").Scan(&paginatio.Total)

	var blogs []Blogs

	vars := mux.Vars(req)
	paginatio.Limi = 3
	paginatio.OFFSE = 0
	paginatio.Last = int(math.Ceil(float64(paginatio.Total) / float64(paginatio.Limi)))

	i, err := strconv.Atoi(vars["pageid"])

	if err != nil {
		paginatio.OFFSE = 1
		log.Println("error while getting page id: ", err.Error())
	} else if i > 0 && ((i-1)*paginatio.Limi) < paginatio.Total {
		paginatio.OFFSE = ((i - 1) * paginatio.Limi)
	} else {
		log.Println("404 page")
		return
	}

	if i == 1 {
		paginatio.Next = i + 1
	} else if i > 1 && i != paginatio.Last-1 {
		paginatio.Next = i + 1
		paginatio.Pre = i - 1
	} else if i == paginatio.Last {
		paginatio.Pre = i - 1
		paginatio.Next = 0
	} else if i > 1 {
		paginatio.Pre = i - 1
	}
	if paginatio.Next == 0 {
		log.Println("there is no Next page ")
	}
	if paginatio.Pre == 0 {
		log.Println("there is no previous page ")
	}
	//log.Println(paginatio.OFFSE)
	//log.Println("page id :",vars["pageid"])

	query := `SELECT  id, created_at, updated_at, message,title FROM blogss OFFSET $1 Limit $2`
	if err := DB.Select(&blogs, query, paginatio.OFFSE, paginatio.Limi); err != nil {
		if err != nil {
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
		return
	}
	var category1 []Category1

	quer := `SELECT * FROM categories`
	if err := DB.Select(&category1, quer); err != nil {
		if err != nil {
			panic(err)
		}
		log.Println("error while getting user from db: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, BlogsTemplate{
		Bloglist:          blogs,
		Paginatio:         paginatio,
		CartegoryDatalist: category1,
	}); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
	//http.Redirect(res, req, "/blogs", http.StatusSeeOther)
}
