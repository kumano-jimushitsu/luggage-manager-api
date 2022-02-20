package models

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type ObjectType interface {
	GetName() string
	Uid() string
}

func ParseJsonToObjects(raw_json string, objectType ObjectType) (interface{}, error) {
	switch objectType.(type) {
	case Ryosei:
		return ParseJsonToRyoseis(raw_json)
	case Parcel:
		return ParseJsonToParcels(raw_json)
	case ParcelEvent:
		return ParseJsonToParcelEvent(raw_json)
	default:
		return nil, errors.New("Unknown object type")
	}
}

func GetAllRecords(db *sqlx.DB, objectType ObjectType) (interface{}, error) {
	switch objectType.(type) {
	case Ryosei:
		return GetAllRyoseis(db)
	case Parcel:
		return GetAllParcels(db)
	case ParcelEvent:
		return GetAllParcelEvents(db)
	default:
		return nil, errors.New("Unknown object type")
	}
}

func InsertObjects(db *sqlx.DB, objects interface{}) error {
	switch objects.(type) {
	case []*Ryosei:
		return InsertRyoseis(db, objects.([]*Ryosei))
	case []*Parcel:
		return InsertParcels(db, objects.([]*Parcel))
	case []*ParcelEvent:
		return InsertParcelEvents(db, objects.([]*ParcelEvent))
	default:
		return errors.New("Unknown objects type")
	}
}

/*
func UpdateObjects(db *sqlx.DB, objects interface{}) error {
	switch objects.(type) {
	case []*Ryosei:
		return UpdateRyoseis(db, objects.([]*Ryosei))
	case []*Parcel:
		return UpdateParcels(db, objects.([]*Parcel))
	default:
		return errors.New("Unknown objects type")
	}
}
*/

func GetUnsyncedObjectsAsSqlInsert(db *sqlx.DB, objectType ObjectType) (*string, error) {
	switch objectType.(type) {
	case Ryosei:
		return GetUnsyncedRyoseisAsSqlInsert(db)
	case Parcel:
		return GetUnsyncedParcelsAsSqlInsert(db)
	case ParcelEvent:
		return GetUnsyncedParcelEventsAsSqlInsert(db)
	default:
		return nil, errors.New("Unknown objects type")
	}
}

/*
func GetUnsyncedObjectsAsSqlUpdate(db *sqlx.DB, objectType ObjectType) (*string, error) {
	switch objectType.(type) {
	case Ryosei:
		return GetUnsyncedRyoseisAsSqlUpdate(db)
	case Parcel:
		return GetUnsyncedParcelsAsSqlUpdate(db)
	default:
		return nil, errors.New("Unknown objects type")
	}
}*/
