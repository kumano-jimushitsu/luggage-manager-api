package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"luggage-api/server/models"
	"net/http"
	"strings"
)

// Route contains information for handlers to run
// The information will be passed to handlers when triggered
type Routes struct {
	rootDir     string
	disableCORS bool
	apiKey      string
}

func (routes *Routes) ryoseiHandler(env *Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/ryosei/")

		if method == "" {
			showAllRyoseis(w, r, env.DB)
		} else if method == "create" {
			processCreateRequest(w, r, env.DB)
		} else if method == "update" {
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[len("/ryosei/"):])
		} else {
			fmt.Fprintf(w, "Wrong path: %s", r.URL.Path[len("/ryosei/"):])
		}
	})
}

/*
	タブレット側からjsonを取り出し、
	jsonをもとにDBを更新する（create）
	最後に成功・失敗ステータスをメッセージで送り返す
*/
func processCreateRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseForm()
	// log.Println(r.Form)
	data_json := r.Form[""]
	sqlcommand := json_test(data_json[0])
	log.Println(sqlcommand)
	fmt.Fprintf(w, "%s", sqlcommand)
}

func showAllRyoseis(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Get all ryoseis
	ryoseis, err := models.AllRyoseis(db)
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
	タブレットでデータベース登録したとき、httpリクエストでタブレット内のデータテーブルがstringで送られる
	stringをjsonにパースして、データベースにinsertする
	更新されたデータベースをタブレット側に送り返す（json）
*/
func (routes *Routes) insertFromTabletHandler(env *Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the string inside http request

		// Insert into the table

		// Send the table to the tablet
	})
}

func (routes *Routes) parcelHandler(env *Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", r.URL.Path[len("/parcel/"):])
	})
}

func json_test(data_json string) interface{} {
	var DB_data []map[string]interface{}
	json.Unmarshal([]byte(data_json), &DB_data)
	result_data := DB_data[0]["ryosei_name"]
	return result_data
}

