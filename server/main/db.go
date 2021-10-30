package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
)

type Env struct {
	DB *sql.DB
}

func NewDB() (*sql.DB, error) {
	query := makeQuery().String()
	db, err := sql.Open("sqlserver", query)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func makeQuery() *url.URL {
	query := url.Values{}
	query.Add("app name", "LuggageAPI")
	query.Add("database", "parcels")
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword("sa", "kumano"),
		Host:     fmt.Sprintf("%s:%d", "127.0.0.1", 1433),
		RawQuery: query.Encode(),
	}
	return u
}

// Connect to MySQL
// db, err := sql.Open("mysql", "root:my-secret-pwd@tcp(db:3306)/mydb")
