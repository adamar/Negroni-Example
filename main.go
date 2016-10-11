package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {

	db, err := setupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	n := negroni.Classic()

	store := cookiestore.New([]byte("ohhhsooosecret"))
	n.Use(sessions.Sessions("global_session_store", store))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		SimplePage(w, r, "mainpage")
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			SimplePage(w, r, "login")
		} else if r.Method == "POST" {
			LoginPost(w, r)
		}
	})

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			SimplePage(w, r, "signup")
		} else if r.Method == "POST" {
			SignupPost(w, r)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		Logout(w, r)
	})

	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		SimpleAuthenticatedPage(w, r, "home")
	})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		APIHandler(w, r)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	n.UseHandler(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3300"
	}
	n.Run(":" + port)
}

func setupDB() (*sql.DB, error) {

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "user=negroni password=negroni dbname=negroni-sample sslmode=disable"
	}
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}

func SimplePage(w http.ResponseWriter, req *http.Request, template string) {
	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func SimpleAuthenticatedPage(w http.ResponseWriter, req *http.Request, template string) {
	session := sessions.GetSession(req)
	sess := session.Get("useremail")

	if sess == nil {
		http.Redirect(w, req, "/notauthenticated", 301)
	}

	r := render.New(render.Options{})
	r.HTML(w, http.StatusOK, template, nil)
}

func LoginPost(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)

	username := req.FormValue("inputUsername")
	password := req.FormValue("inputPassword")

	var (
		email                string
		password_in_database string
	)

	err := db.QueryRow("SELECT user_email, user_password FROM users WHERE user_name = $1", username).Scan(&email, &password_in_database)
	if err == sql.ErrNoRows {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_in_database), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		http.Redirect(w, req, "/authfail", 301)
	} else if err != nil {
		log.Print(err)
		http.Redirect(w, req, "/authfail", 301)
	}

	session.Set("useremail", email)
	http.Redirect(w, req, "/home", 302)
}

func SignupPost(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("inputUsername")
	password := req.FormValue("inputPassword")
	email := req.FormValue("inputEmail")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO users (user_name, user_password, user_email) VALUES ($1, $2, $3)", username, string(hashedPassword), email)
	if err != nil {
		log.Print(err)
	}

	http.Redirect(w, req, "/login", 302)
}

func Logout(w http.ResponseWriter, req *http.Request) {
	session := sessions.GetSession(req)
	session.Delete("useremail")
	http.Redirect(w, req, "/", 302)
}

func APIHandler(w http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal("{'API Test':'Works!'}")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
