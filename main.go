package main

import (
	"database/sql"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	_ "github.com/lib/pq"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"os"
)

var db *sql.DB = setupDB()

func init() {

	db.Exec(`CREATE TABLE users (
                 id SERIAL,
                 user_name VARCHAR(60),  
                 user_email VARCHAR(60),  
                 user_password VARCHAR(60),  
                 user_created TIMESTAMP WITH TIME ZONE,
                 user_last_login TIMESTAMP WITH TIME ZONE, 
                 PRIMARY KEY  (id),  
                 CONSTRAINT users_email UNIQUE (user_email)
            );`)

	db.Exec(`INSERT INTO users (user_name, user_email, user_password)
             VALUES ('john', 'john@example.com', 'supersecret');`)

}

func main() {

	defer db.Close()

	mux := http.NewServeMux()
	n := negroni.Classic()

	store := sessions.NewCookieStore([]byte("ohhhsooosecret"))
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
	n.Run(":" + os.Getenv("PORT"))

}

func setupDB() *sql.DB {

	db_url := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}

	return db

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
		email string
	)

    err := db.QueryRow("SELECT user_email FROM users WHERE user_name = $1 AND user_password = $2", username, password).Scan(&email)
	if err != nil {
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

	_, err := db.Exec("INSERT INTO users (user_name, user_password, user_email) VALUES ($1, $2, $3)", username, password, email)
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
