package models

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type ParcelEvent struct {
	ObjectType         ObjectType
	Id                 string         `json:"uid" db:"uid"`
	CreatedAt          string         `json:"created_at" db:"created_at"`
	EventType          int            `json:"event_type" db:"event_type"`
	ParcelUid          string         `json:"parcel_uid" db:"parcels_uid"`
	RyoseiUid          string         `json:"ryosei_uid" db:"ryosei_uid"`
	RoomID             string         `json:"room_name" db:"room_name"`
	Name               string         `json:"ryosei_name" db:"ryosei_name"`
	TargetID           sql.NullString `json:"target_event_uid" db:"target_event_uid"`
	Note               sql.NullString `json:"note" db:"note"`
	AfterPeriodicCheck int            `json:"is_after_fixed_time" db:"is_after_fixed_time"`
	IsFinished         int            `json:"is_finished" db:"is_finished"`
	IsDeleted          int            `json:"is_deleted" db:"is_deleted"`
	SharingStatus      int            `json:"sharing_status" db:"sharing_status"`
}

func (parcelEvent ParcelEvent) GetName() string {
	return "parcelEvent"
}

func ParseJsonToParcelEvent(raw_json string) ([]*ParcelEvent, error) {
	var events []*ParcelEvent
	err := json.Unmarshal([]byte(raw_json), &events)
	if err != nil {
		return nil, err
	}
	return events, err
}

func getEventFromSqlRows(db *sqlx.DB) ([]*ParcelEvent, error) {
	rows, err := db.Query("SELECT * FROM parcel_event WHERE sharing_status = 20 OR sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*ParcelEvent, 0)
	for rows.Next() {
		event := new(ParcelEvent)
		err := rows.Scan(
			&event.Id,
			&event.CreatedAt,
			&event.EventType,
			&event.ParcelUid,
			&event.RyoseiUid,
			&event.RoomID,
			&event.Name,
			&event.TargetID,
			&event.Note,
			&event.AfterPeriodicCheck,
			&event.IsFinished,
			&event.IsDeleted,
			&event.SharingStatus,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func GetAllEvent(db *sqlx.DB) ([]*ParcelEvent, error) {
	rows, err := db.Query("SELECT * FROM parcel_event")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*ParcelEvent, 0)
	for rows.Next() {
		event := new(ParcelEvent)
		err := rows.Scan(
			&event.Id,
			&event.CreatedAt,
			&event.EventType,
			&event.ParcelUid,
			&event.RyoseiUid,
			&event.RoomID,
			&event.Name,
			&event.TargetID,
			&event.Note,
			&event.AfterPeriodicCheck,
			&event.IsFinished,
			&event.IsDeleted,
			&event.SharingStatus,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func setEvent(event *ParcelEvent, record *map[string]interface{}) error {
	event.Id = (*record)["uid"].(string)
	event.CreatedAt = (*record)["created_at"].(string)
	event.EventType = floatToInt((*record)["event_type"].(float64))
	event.ParcelUid = (*record)["parcel_uid"].(string)
	event.RyoseiUid = (*record)["ryosei_uid"].(string)
	event.RoomID = (*record)["room_name"].(string)
	event.Name = (*record)["ryosei_name"].(string)
	event.TargetID = toNullString((*record)["target_event_uid"].(string))
	event.Note = toNullString((*record)["note"].(string))
	event.AfterPeriodicCheck = floatToInt((*record)["is_after_fixed_time"].(float64))
	event.IsFinished = floatToInt((*record)["is_finished"].(float64))
	event.IsDeleted = floatToInt((*record)["is_deleted"].(float64))
	event.SharingStatus = floatToInt((*record)["sharing_status"].(float64))
	return nil
}

func GetAllParcelEvents(db *sqlx.DB) ([]*ParcelEvent, error) {
	return []*ParcelEvent{}, nil
}

func InsertParcelEvents(db *sqlx.DB, events []*ParcelEvent) error {
	return nil
}

func UpdateParcelEvents(db *sqlx.DB, events []*ParcelEvent) error {
	return nil
}

func GetUnsyncedParcelEventsAsSqlInsert(db *sqlx.DB) (*string, error) {
	sql := ""
	return &sql, nil
}

func GetUnsyncedParcelEventsAsSqlUpdate(db *sqlx.DB) (*string, error) {
	sql := ""
	return &sql, nil
}
