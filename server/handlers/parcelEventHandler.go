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

func parcelEventHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/parcelEvent/")

		if method == "" {
			showAllParcelEvents(w, r, env.DB)
		} else if method == "create" {
			createParcelEvent(w, r, env.DB)
		} else if method == "update" {
			updateParcelEvent(w, r, env.DB)
		} else if method == "check" {
			checkParcelEventUpdateInTablet(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/parcelEvent/"):])
		}
	})
}

func showAllParcelEvents(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

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

func createParcelEvent(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	ryoseis, err := parseParcelEventJson(raw_json)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", ryoseis[0].Id)
}

func updateParcelEvent(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	raw_json := r.Form[""][0]
	ryoseis, err := parseRyoseiJson(raw_json)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", ryoseis[0].Id)
}

func parseParcelEventJson(raw_json string) ([]*models.Ryosei, error) {
	var ryoseis []*models.Ryosei
	err := json.Unmarshal([]byte(raw_json), &ryoseis)
	if err != nil {
		return nil, err
	}
	return ryoseis, err
}

// 20 or 21 -> 30に更新
func checkParcelEventUpdateInTablet(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	r.ParseForm()
	msg := r.Form[""][0]
	if msg == "Success" {
		update := "UPDATE parcels SET sharing_status = 30 WHERE sharing_status == 20 OR sharing_status == 21"
		_, err := db.Exec(update)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// do nothing
	}
	fmt.Fprintf(w, "%s", "")
}
