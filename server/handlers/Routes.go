package handlers

import (
	"fmt"
	"luggage-api/server/database"
	"luggage-api/server/models"
	"net/http"
)

// Route contains information for handlers to run
// The information will be passed to handlers when triggered
type Routes struct {
	RootDir     string
	DisableCORS bool
	ApiKey      string
}

func (routes *Routes) ObjectHandler(env *database.Env, objectType models.ObjectType) http.Handler {
	switch t := objectType.(type) {
	case models.Ryosei:
		return ryoseiHandler(env)
	case models.Parcel:
		return parcelHandler(env)
	case models.ParcelEvent:
		return parcelEventHandler(env)
	default:
		return exceptionHandler(env, t)
	}
}

func exceptionHandler(env *database.Env, t models.ObjectType) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Unknown Type %T", t)
	})
}
