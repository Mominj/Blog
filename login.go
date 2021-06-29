package main

import (
	"html/template"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

type LoginTemplate struct {
	LoginForm
	Errors map[string]error
}

func (l LoginForm) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email,
			validation.Required.Error("Email is required"),
		),
		validation.Field(&l.Password,
			validation.Required.Error("Password is required"),
			validation.Length(6, 16).Error("Password must be 6 to 16 characters length"),
		),
	)
}

func loginForm(res http.ResponseWriter, req *http.Request) {
	viewLoginForm(res, LoginTemplate{})
}
func init() {
	DB, DBErr = sqlx.Connect("postgres", "user=postgres password=password dbname=new sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while connecting to database", DBErr)
	}

	SessionStore = sessions.NewCookieStore([]byte("SECRET-KEY"))

}
func login(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error while parsing form: ", err.Error())
	}

	var loginForm LoginForm

	err = decoder.Decode(&loginForm, req.PostForm)
	if err != nil {
		log.Println("error while decoding form to struct: ", err.Error())
	}

	if vErr := loginForm.Validate(); vErr != nil {
		log.Println("failed to validate form: ", loginForm)
		if vErrs, ok := vErr.(validation.Errors); ok {
			viewLoginForm(res, LoginTemplate{
				LoginForm: loginForm,
				Errors:    vErrs,
			})
			return
		}
	}

	var user User

	query := `SELECT  id, email, password FROM userss where email = $1`
	if err := DB.Get(&user, query, loginForm.Email); err != nil {
		if err != nil {
			panic(err)
		}

		log.Println("error while getting user from db: ", err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password)); err != nil {
		panic(err)
	}

	session, err := SessionStore.Get(req, "first_app-session")
	if err != nil {
		log.Println("error while getting session: ", err.Error())
		return
	}

	session.Values["user_id"] = user.ID

	if err := session.Save(req, res); err != nil {
		log.Println("error while saving session: ", err.Error())
		return
	}

	http.Redirect(res, req, "/blogsdetails/1", http.StatusSeeOther)
}

func viewLoginForm(res http.ResponseWriter, loginTemplate LoginTemplate) {
	parsedTemplate, err := template.ParseFiles("templates/auth/login.html")
	if err != nil {
		log.Println("error while parsing template: ", err.Error())
		return
	}

	if err := parsedTemplate.Execute(res, loginTemplate); err != nil {
		log.Println("error while executing template: ", err.Error())
	}
}
