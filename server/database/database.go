package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
)

type Env struct {
	DB *sqlx.DB
}

func NewDB(dbName string) (*sqlx.DB, error) {
	query := makeQuery(dbName).String()
	// Connect to MySQL
	// db, err := sql.Open("mysql", "root:my-secret-pwd@tcp(db:3306)/mydb")
	// Connect to MSSQL

	db, err := sqlx.Open("sqlserver", query)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db, err
}

func GetTestTransaction(dbName string) *sql.Tx {
	query := makeQuery(dbName).String()
	db, err := sqlx.Open("sqlserver", query)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tx, _ := db.Begin()

	return tx
}

func makeQuery(dbName string) *url.URL {
	query := url.Values{}
	query.Add("app name", "LuggageAPI")
	query.Add("database", dbName)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword("sa", "kumano"),
		Host:     fmt.Sprintf("%s:%d", "127.0.0.1", 1433),
		RawQuery: query.Encode(),
	}
	return u
}
