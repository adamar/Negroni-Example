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

var db *sql.DB = setupDB()

func main() {

	defer db.Close()

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

func setupDB() *sql.DB {
	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "user=negroni password=negroni dbname=negroni-sample sslmode=disable"
	}
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return db
}

func errHandler(err error) {
	if err != nil {
		log.Print(err)
	}
}
