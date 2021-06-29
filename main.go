package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

var decoder = schema.NewDecoder()

var DB *sqlx.DB
var DBErr error

var SessionStore *sessions.CookieStore

func home(res http.ResponseWriter, req *http.Request) {

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

	parsedTemplate, err := template.New("home.html").Funcs(funcMap).ParseFiles("templates/home.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, nil); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
}

type User struct {
	ID        int32  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
}

type UserCreateTemplate struct {
	User
	Errors map[string]error
}

func userCreate(res http.ResponseWriter, req *http.Request) {
	ViewUserCreateTemplate(res, UserCreateTemplate{})
}
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName,
			validation.Required.Error("First name is required"),
			validation.Length(2, 10).Error("First name must be 3 to 10 characters length"),
		),
		validation.Field(&u.LastName,
			validation.Required.Error("Last name is required"),
			validation.Length(2, 10).Error("Last name must be 3 to 10 characters length"),
		),
		validation.Field(&u.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Must be a valid email"),
		),
		validation.Field(&u.Password,
			validation.Required.Error("Password is required"),
			validation.Length(6, 16).Error("Password must be 6 to 16 characters length"),
		),
	)
}
func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=password dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}

	SessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

}
func userStore(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error while parsing form: ", err.Error())
	}

	var user User

	err = decoder.Decode(&user, req.PostForm)
	if err != nil {
		log.Println("error while decoding form to struct: ", err.Error())
	}

	if vErr := user.Validate(); vErr != nil {
		log.Println("failed to validate form: ", user)
		if vErrs, ok := vErr.(validation.Errors); ok {
			fmt.Println(vErrs)
			viewUserCreateTemplate := UserCreateTemplate{
				User:   user,
				Errors: vErrs,
			}
			ViewUserCreateTemplate(res, viewUserCreateTemplate)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while encrypted password: ", err.Error())
	}
	user.Password = string(hash)

	query := `
		INSERT INTO userss(
			first_name,
			last_name,
			email,
			password
		)
		VALUES(
			:first_name,
			:last_name,
			:email,
			:password
		)
		RETURNING id
	`

	var id int32
	stmt, err := DB.PrepareNamed(query)

	if err != nil {
		log.Println("db error: failed prepare ", err.Error())
		return
	}

	if err := stmt.Get(&id, user); err != nil {
		log.Println("db error: failed to insert data ", err.Error())
		return
	}

	http.Redirect(res, req, "/login", http.StatusSeeOther)
}

func ViewUserCreateTemplate(res http.ResponseWriter, userCreateTemplate UserCreateTemplate) {
	parsedTemplate, err := template.ParseFiles("templates/user/user.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, userCreateTemplate); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
}
func main() {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	router.HandleFunc("/home", home).Methods("GET").Name("home")
	router.HandleFunc("/users/create", userCreate).Methods("GET")
	router.HandleFunc("/users", userStore).Methods("POST")
	router.HandleFunc("/login", loginForm).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("GET").Name("logout")

	router.HandleFunc("/createblog", createblog).Methods("GET")
	router.HandleFunc("/createblog", blogStore).Methods("POST")
	router.HandleFunc("/blogsdetails", getblogs).Methods("GET")
	router.HandleFunc("/blogsdetails/{pageid}", getblogs).Methods("GET")

	router.HandleFunc("/blogscategorycreate", blogsCategoryCreate).Methods("GET")
	router.HandleFunc("/blogscategorycreate", blogsCategoryStore).Methods("POST")
	router.HandleFunc("/blogscategory/{category_id}", blogsCategoryShow).Methods("GET")
	router.HandleFunc("/blogscat/{category_id}", blogsCategoryShow1).Methods("GET")
	router.HandleFunc("/blogscat/category_id/{pageid}", blogsCategoryShow1).Methods("GET")

	router.HandleFunc("/blogs/{id}", singleBlogShow).Methods("GET")
	router.HandleFunc("/blogs/{id}/edit", blogEdit).Methods("GET")
	router.HandleFunc("/blogs/{id}/edit", blogUpdate).Methods("POST")
	router.HandleFunc("/blogs/{id}/delete", blogDelete).Methods("GET")
	router.HandleFunc("/blogs/{id}", commnetStore).Methods("POST")

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", CONN_HOST, CONN_PORT), router); err != nil {
		log.Fatal("error starting server: ", err)
	}
}
