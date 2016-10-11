package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {

	db, err := setupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
                 id SERIAL,
                 user_name VARCHAR(60),  
                 user_email VARCHAR(60),  
                 user_password VARCHAR(60),  
                 user_created TIMESTAMP WITH TIME ZONE,
                 user_last_login TIMESTAMP WITH TIME ZONE, 
                 PRIMARY KEY  (id),  
                 CONSTRAINT users_email UNIQUE (user_email)
            );`)

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
