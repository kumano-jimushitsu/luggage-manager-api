package main

import (
	"log"
	"luggage-api/server/database"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Connect to Database
	db, err := database.NewDB("parcels")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// env = a pacckage of Global environment varible
	env := &database.Env{DB: db}

	// Define a non-default ServeMux
	mux := http.NewServeMux()

	routes := Routes{
		rootDir:     "",
		disableCORS: true,
		apiKey:      "aaa",
	}

	// Register event handlers
	mux.Handle("/ryosei/", routes.ryoseiHandler(env))
	mux.Handle("/parcel/", routes.parcelHandler(env))

	// Start the Server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("Can't start server. Shutting down...")
	}
}
