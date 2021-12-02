package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"luggage-api/server/models"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func ShowAllParcels(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

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

func CreateParcel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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

func UpdateParcel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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

func CheckParcelUpdateInTablet(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
