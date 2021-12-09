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

func parcelHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/parcel/")

		if method == "" {
			showAllParcels(w, r, env.DB)
		} else if method == "create" {
			createParcel(w, r, env.DB)
		} else if method == "update" {
			updateParcel(w, r, env.DB)
		} else if method == "check" {
			checkParcelUpdateInTablet(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/parcel/"):])
		}
	})
}

func showAllParcels(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	// Get all ryoseis
	parcels, err := models.GetAllParcels(db)
	if err != nil {
		log.Fatal(err)
	}

	// process to json
	json, err := json.Marshal(parcels)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s", string(json))
}

func createParcel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	if raw_json != "" {
		parcels, err := parseParcelJson(raw_json)
		if err != nil {
			log.Fatal(err)
		}
		err = models.InsertParcels(db, parcels)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := models.GetUnsyncedParcelsAsSqlInsert(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}

func updateParcel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	if raw_json != "" {
		parcels, err := parseParcelJson(raw_json)
		if err != nil {
			log.Fatal(err)
		}
		err = models.UpdateParcels(db, parcels)
		if err != nil {
			log.Fatal(err)
		}
	}
	msg, err := models.GetUnsyncedParcelsAsSqlInsert(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}

func parseParcelJson(raw_json string) ([]*models.Parcel, error) {
	var parcels []*models.Parcel
	err := json.Unmarshal([]byte(raw_json), &parcels)
	if err != nil {
		return nil, err
	}
	return parcels, err
}

func checkParcelUpdateInTablet(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	msg := r.Form[""][0]
	if msg == "Success" {
		// 20 or 21 -> 30に更新
		update := "UPDATE parcels SET sharing_status = 30 WHERE sharing_status = 20 OR sharing_status = 21"
		_, err := db.Exec(update)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// do nothing
	}
	fmt.Fprintf(w, "%s", "")
}
