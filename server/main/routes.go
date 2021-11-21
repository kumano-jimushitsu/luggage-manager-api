package main

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
	rootDir     string
	disableCORS bool
	apiKey      string
}

func (routes *Routes) ryoseiHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Trim the action from the request url
		method := strings.TrimPrefix(r.URL.Path, "/ryosei/")

		if method == "" {
			showAllRyoseis(w, r, env.DB)
		} else if method == "create" {
			createRyosei(w, r, env.DB)
		} else if method == "update" {
			updateRyosei(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/ryosei/"):])
		}
	})
}

func (routes *Routes) parcelHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/parcel/")

		if method == "" {
			showAllParcels(w, r, env.DB)
		} else if method == "create" {
			createParcel(w, r, env.DB)
		} else if method == "update" {
			updateParcel(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/parcel/"):])
		}
	})
}

func showAllRyoseis(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	// Get all ryoseis
	ryoseis, err := models.GetAllRyoseis(db)
	if err != nil {
		log.Fatal(err)
	}

	// process to json
	json, err := json.Marshal(ryoseis)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s", string(json))
}

/*
	タブレット側からjsonを取り出し、
	jsonをもとにDBを更新する（create）
	最後に成功・失敗ステータスをメッセージで送り返す
*/
func createRyosei(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	ryoseis, err := parseRyoseiJson(raw_json)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", ryoseis[0].Id)
}

func updateRyosei(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

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
	parcels, err := parseParcelJson(raw_json)
	if err != nil {
		log.Fatal(err)
	}
	err = models.InsertParcels(db, parcels)
	if err != nil {
		log.Fatal(err)
	}
	msg, err := models.GetUnsyncedParcelsAsSqlInsert(db)
	log.Fatal(*msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", *msg)
}

func updateParcel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	fmt.Println("aaa")
}

func parseRyoseiJson(raw_json string) ([]*models.Ryosei, error) {
	var ryoseis []*models.Ryosei
	err := json.Unmarshal([]byte(raw_json), &ryoseis)
	if err != nil {
		return nil, err
	}
	return ryoseis, err
}

func parseParcelJson(raw_json string) ([]*models.Parcel, error) {
	var parcels []*models.Parcel
	err := json.Unmarshal([]byte(raw_json), &parcels)
	if err != nil {
		return nil, err
	}
	return parcels, err
}
