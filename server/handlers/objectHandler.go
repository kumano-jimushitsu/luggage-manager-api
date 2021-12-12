package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"luggage-api/server/database"
	"luggage-api/server/models"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Route contains information for handlers to run
// The information will be passed to handlers when triggered
type Routes struct {
	RootDir     string
	DisableCORS bool
	ApiKey      string
}

func (routes *Routes) ObjectHandler(env *database.Env, objectType models.ObjectType) http.Handler {
	prefix := "/" + objectType.GetName() + "/"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, prefix)
		switch method {
		case "":
			showAllObjects(w, r, env.DB, objectType)
		case "create":
			createObjects(w, r, env.DB, objectType)
		case "update":
			updateObjects(w, r, env.DB, objectType)
		case "check":
			checkObjectUpdateInTablet(w, r, env.DB, objectType)
		default:
			fmt.Fprintf(w, "Wrong action: %v", r.URL.Path[len(prefix)])
		}
	})
}

func showAllObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objectType models.ObjectType) {

	// Get all objects from database
	objects, err := models.GetAllRecords(db, objectType)
	if err != nil {
		log.Fatal(err)
	}

	// Process objets to json
	json, err := json.Marshal(objects)
	if err != nil {
		log.Fatal(err)
	}

	// Output
	fmt.Fprintf(w, "%s", string(json))
}

func createObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objectType models.ObjectType) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	if raw_json != "" {
		objects, err := models.ParseJsonToObjects(raw_json, objectType)
		if err != nil {
			log.Fatal(err)
		}
		err = models.InsertObjects(db, objects)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := models.GetUnsyncedObjectsAsSqlInsert(db, objectType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}

func updateObjects(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objeType models.ObjectType) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	if raw_json != "" {
		objects, err := models.ParseJsonToObjects(raw_json, objeType)
		if err != nil {
			log.Fatal(err)
		}
		err = models.UpdateObjects(db, objects)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := models.GetUnsyncedObjectsAsSqlUpdate(db, objeType)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}

func checkObjectUpdateInTablet(w http.ResponseWriter, r *http.Request, db *sqlx.DB, objectType models.ObjectType) {

	r.ParseForm()
	msg := r.Form[""][0]

	var update string
	switch objectType.(type) {
	case models.Ryosei:
		update = "UPDATE ryosei SET sharing_status = 30 WHERE sharing_status = 20 OR sharing_status = 21"
	case models.Parcel:
		update = "UPDATE parcels SET sharing_status = 30 WHERE sharing_status = 20 OR sharing_status = 21"
	case models.ParcelEvent:
		update = "UPDATE parcel_event SET sharing_status = 30 WHERE sharing_status = 20 OR sharing_status = 21"
	default:
		log.Fatal("Unknown type")
	}

	if msg == "Success" {
		_, err := db.Exec(update)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// do nothing
	}
	fmt.Fprintf(w, "%s", "")
}
