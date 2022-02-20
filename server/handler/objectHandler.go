package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"luggage-api/server/database"
	"luggage-api/server/models"
	"net/http"
	"strings"
	"time"
)

// Route contains information for handlers to run
// The information will be passed to handlers when triggered
type Routes struct {
	RootDir string
	// DisableCORS bool
	// ApiKey      string
}

func (routes *Routes) ObjectHandler(env *database.Env, objectType models.ObjectType) http.Handler {
	prefix := "/" + objectType.GetName() + "/"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, prefix)
		switch method {
		case "":
			showAllObjects(w, r, env, objectType)
		case "create":
			createObjects(w, r, env, objectType)
		//case "update":
		// updateObjects(w, r, env.DB, objectType)
		// case "check":
		// checkObjectUpdateInTablet(w, r, env.DB, objectType)
		default:
			fmt.Fprintf(w, "Wrong action: %v", r.URL.Path[len(prefix)])
		}
	})
}

func showAllObjects(w http.ResponseWriter, r *http.Request, env *database.Env, objectType models.ObjectType) {

	// Get all objects from database
	objects, err := models.GetAllRecords(env.DB, objectType)
	if err != nil {
		log.Fatal(err)
		env.Logger.Fatal(err)
	}

	// Process objets to json
	json, err := json.Marshal(objects)
	if err != nil {
		log.Fatal(err)
		env.Logger.Fatal(err)
	}

	// Output
	fmt.Fprintf(w, "%s", string(json))
}

func createObjects(w http.ResponseWriter, r *http.Request, env *database.Env, objectType models.ObjectType) {
	r.ParseForm()
	raw_json := r.Form[""][0]

	if raw_json != "" {

		// Objectify json
		objects, err := models.ParseJsonToObjects(raw_json, objectType)
		if err != nil {
			log.Fatal(err)
			env.Logger.Fatal(err)
		}

		// Log message in console
		consoleLog(env, objects)

		// Upsert objects
		err = models.InsertObjects(env.DB, objects)
		if err != nil {
			log.Fatal(err)
			env.Logger.Fatal(err)
		}

	}

	// Send objects with sharing_status = 20 to the tablet
	msg, err := models.GetUnsyncedObjectsAsSqlInsert(env.DB, objectType)
	if err != nil {
		log.Fatal(err)
		env.Logger.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)

}

// 処理が汚い。[]*any見たいな、中身の型はわからないけどこれは配列型だよっていう表現をできたらいいのだが
func consoleLog(env *database.Env, objects interface{}) {
	switch objects := objects.(type) {
	case []*models.Ryosei:
		for _, object := range objects {
			currentTime := time.Now().Format("2006-01-02-15:04:05")
			fmt.Printf("[%v] Received %v upsert data with uid = %v\n", currentTime, (*object).GetName(), (*object).Uid())
			env.Logger.Printf("\"Received %v upsert data with uid = %v\"\n", (*object).GetName(), (*object).Uid())
		}
	case []*models.Parcel:
		for _, object := range objects {
			currentTime := time.Now().Format("2006-01-02-15:04:05")
			fmt.Printf("[%v] Received %v upsert data with uid = %v\n", currentTime, (*object).GetName(), (*object).Uid())
			env.Logger.Printf("\"Received %v upsert data with uid = %v\"\n", (*object).GetName(), (*object).Uid())
		}
	case []*models.ParcelEvent:
		for _, object := range objects {
			currentTime := time.Now().Format("2006-01-02-15:04:05")
			fmt.Printf("[%v] Received %v upsert data with uid = %v\n", currentTime, (*object).GetName(), (*object).Uid())
			env.Logger.Printf("\"Received %v upsert data with uid = %v\"\n", (*object).GetName(), (*object).Uid())
		}
	default:
		fmt.Println("Unknown object type")
		env.Logger.Println("Unknown object type")
	}
}

/*
func updateObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objectType models.ObjectType) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	if raw_json != "" {
		objects, err := models.ParseJsonToObjects(raw_json, objectType)
		if err != nil {
			log.Fatal(err)
		}
		err = models.UpdateObjects(db, objects)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := models.GetUnsyncedObjectsAsSqlUpdate(db, objectType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}
*/

/* Check if sql insert in table succeeded */
// func checkObjectUpdateInTablet(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objectType models.ObjectType) {

// 	r.ParseForm()
// 	msg := r.Form[""][0]
// 	msg = msg[:6]

// 	if msg == "" {
// 		// do nothing
// 		return
// 	}

// 	var sharing_status int

// 	switch msg {
// 	case "create":
// 		sharing_status = 20
// 	case "update":
// 		sharing_status = 21
// 	default:
// 		log.Fatal("Unknown method")
// 	}

// 	var update string

// 	switch objectType.(type) {
// 	case models.Ryosei:
// 		update = "UPDATE ryosei SET sharing_status = 30 WHERE sharing_status = " + fmt.Sprint(sharing_status)
// 	case models.Parcel:
// 		update = "UPDATE parcels SET sharing_status = 30 WHERE sharing_status = " + fmt.Sprint(sharing_status)
// 	case models.ParcelEvent:
// 		update = "UPDATE parcel_event SET sharing_status = 30 WHERE sharing_status = " + fmt.Sprint(sharing_status)
// 	default:
// 		log.Fatal("Unknown type")
// 	}

// 	_, err := db.Exec(update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Fprintf(w, "%s", "")
// }

// /* Bulk insert ryosei from csv */
// func (routes *Routes) InitRyoseiHandler(env *database.Env) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		models.GetRyoseiSeedingCsv(env.DB)
// 	})
// }
