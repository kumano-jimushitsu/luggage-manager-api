package main

import (
	"fmt"
	"luggage-api/server/database"
	"luggage-api/server/handlers"
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

func (routes *Routes) ryoseiHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Trim the action from the request url
		method := strings.TrimPrefix(r.URL.Path, "/ryosei/")

		if method == "" {
			handlers.ShowAllRyoseis(w, r, env.DB)
		} else if method == "create" {
			handlers.CreateRyosei(w, r, env.DB)
		} else if method == "update" {
			handlers.UpdateRyosei(w, r, env.DB)
		} else if method == "check" {
			handlers.CheckRyoseiUpdateInTablet(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/ryosei/"):])
		}
	})
}

func (routes *Routes) parcelHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/parcel/")

		if method == "" {
			handlers.ShowAllParcels(w, r, env.DB)
		} else if method == "create" {
			handlers.CreateParcel(w, r, env.DB)
		} else if method == "update" {
			handlers.UpdateParcel(w, r, env.DB)
		} else if method == "check" {
			handlers.CheckParcelUpdateInTablet(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/parcel/"):])
		}
	})
}

func (routes *Routes) parcelEventHandler(env *database.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/parcelEvent/")

		if method == "" {
			handlers.ShowAllParcelEvents(w, r, env.DB)
		} else if method == "create" {
			handlers.CreateParcelEvent(w, r, env.DB)
		} else if method == "update" {
			handlers.UpdateParcelEvent(w, r, env.DB)
		} else if method == "check" {
			handlers.CheckParcelEventUpdateInTablet(w, r, env.DB)
		} else {
			fmt.Fprintf(w, "Wrong action: %s", r.URL.Path[len("/parcelEvent/"):])
		}
	})
}
