package main

import (
	"fmt"
	"log"
	"luggage-api/server/database"
	handlers "luggage-api/server/handler"
	"luggage-api/server/models"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Output address of log
	currentTime := time.Now().Format("20060102-150405")
	logPath := fmt.Sprintf("../log-%v.txt", currentTime)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	// Instanciate logger
	logger := log.New(f, "", log.LstdFlags)

	// Connect to Database
	db, err := database.NewDB("parcels")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// env = a pacckage of Global environment varible
	env := &database.Env{DB: db, Logger: logger}

	// Define a non-default ServeMux
	mux := http.NewServeMux()

	routes := handlers.Routes{
		RootDir: "",
		// DisableCORS: true,
		// ApiKey:      "aaa",
	}

	// Register event handlers
	mux.Handle("/ryosei/", routes.ObjectHandler(env, models.Ryosei{}))
	mux.Handle("/parcels/", routes.ObjectHandler(env, models.Parcel{}))
	mux.Handle("/parcel_event/", routes.ObjectHandler(env, models.ParcelEvent{}))
	// mux.Handle("/initRyosei", routes.InitRyoseiHandler(env))

	// Start the Server
	log.Println("Server started")
	env.Logger.Println("Server started")
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatalln("Can't start server. Shutting down...")
		env.Logger.Fatalln("Can't start server. Shutting down...")
	}
}
