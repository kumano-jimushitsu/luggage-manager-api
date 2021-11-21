package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"luggage-api/server/database"
	"luggage-api/server/models"
	"os"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestUnmarshalJSON(t *testing.T) {
	tx := database.GetTestTransaction("testDatabase")

	raw_json, err := ioutil.ReadFile("./json/parcel_10.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var parcels []*models.Parcel
	err = json.Unmarshal([]byte(raw_json), &parcels)
	if err != nil {
		log.Fatal(err)
	}

	if parcels[0].OwnerRyoseiName != "三好宏美" {
		t.Errorf("Unmarshal failed. parcel[0].OwnerRyoseiName got %v, want %v", parcels[0].OwnerRyoseiName, "三好宏美")
	}

	tx.Rollback()
}
