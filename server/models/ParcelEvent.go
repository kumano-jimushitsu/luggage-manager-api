package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

/*
	Implement ObjectType interface
*/
func (parcelEvent ParcelEvent) GetName() string {
	return "parcel_event"
}

/*
	Implement Unmarshaler interface
*/
func (parcelEvent *ParcelEvent) UnmarshalJSON(data []byte) error {

	// get parcel events
	var record *map[string]interface{}
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	// parse number(float64) type to bool
	err = setParcelEvent(parcelEvent, record)
	if err != nil {
		return err
	}

	return err
}

func setParcelEvent(parcelEvent *ParcelEvent, record *map[string]interface{}) error {
	parcelEvent.Id = (*record)["uid"].(string)
	parcelEvent.CreatedAt = (*record)["created_at"].(string)
	parcelEvent.EventType = floatToInt((*record)["event_type"].(float64))
	parcelEvent.ParcelUid = toNullString((*record)["parcel_uid"])
	parcelEvent.RyoseiUid = toNullString((*record)["ryosei_uid"])
	parcelEvent.RoomID = toNullString((*record)["room_name"])
	parcelEvent.Name = toNullString((*record)["ryosei_name"])
	parcelEvent.TargetID = toNullString((*record)["target_event_uid"])
	parcelEvent.Note = toNullString((*record)["note"])
	parcelEvent.IsAfterPeriodicCheck = floatToInt((*record)["is_after_fixed_time"].(float64))
	parcelEvent.IsFinished = floatToInt((*record)["is_finished"].(float64))
	parcelEvent.IsDeleted = floatToInt((*record)["is_deleted"].(float64))
	parcelEvent.SharingStatus = floatToInt((*record)["sharing_status"].(float64))

	return nil
}

type ParcelEvent struct {
	Id                   string         `json:"uid" db:"uid"`
	CreatedAt            string         `json:"created_at" db:"created_at"`
	EventType            int            `json:"event_type" db:"event_type"`
	ParcelUid            sql.NullString `json:"parcel_uid" db:"parcel_uid"`
	RyoseiUid            sql.NullString `json:"ryosei_uid" db:"ryosei_uid"`
	RoomID               sql.NullString `json:"room_name" db:"room_name"`
	Name                 sql.NullString `json:"ryosei_name" db:"ryosei_name"`
	TargetID             sql.NullString `json:"target_event_uid" db:"target_event_uid"`
	Note                 sql.NullString `json:"note" db:"note"`
	IsAfterPeriodicCheck int            `json:"is_after_fixed_time" db:"is_after_fixed_time"`
	IsFinished           int            `json:"is_finished" db:"is_finished"`
	IsDeleted            int            `json:"is_deleted" db:"is_deleted"`
	SharingStatus        int            `json:"sharing_status" db:"sharing_status"`
}

/*
	Parse json and return the slice of struct ParcelEvent
*/
func ParseJsonToParcelEvent(raw_json string) ([]*ParcelEvent, error) {
	var events []*ParcelEvent
	err := json.Unmarshal([]byte(raw_json), &events)
	if err != nil {
		return nil, err
	}
	return events, err
}

/*
	Return all records in the parcel_event table
*/
func GetAllParcelEvents(db *sqlx.DB) ([]*ParcelEvent, error) {
	rows, err := db.Query("SELECT * FROM parcel_event")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events, err := getParcelEventsFromSqlRows(db, rows)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func getParcelEventsFromSqlRows(db *sqlx.DB, rows *sql.Rows) ([]*ParcelEvent, error) {
	var isAfterPeriodicCheck bool
	var isFinished bool
	var isDeleted bool
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
			&isAfterPeriodicCheck,
			&isFinished,
			&isDeleted,
			&event.SharingStatus,
		)
		if err != nil {
			return nil, err
		}

		event.IsAfterPeriodicCheck = boolToInt(isAfterPeriodicCheck)
		event.IsFinished = boolToInt(isFinished)
		event.IsDeleted = boolToInt(isDeleted)

		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

var parcelEventInsert string = `
INSERT INTO parcel_event(
	uid,
	created_at,
	event_type,
	parcel_uid,
	ryosei_uid,
	room_name,
	ryosei_name,
	target_event_uid,
	note,
	is_after_fixed_time,
	is_finished,
	is_deleted,
	sharing_status
) VALUES (
	:uid,
	:created_at,
	:event_type,
	:parcel_uid,
	:ryosei_uid,
	:room_name,
	:ryosei_name,
	:target_event_uid,
	:note,
	:is_after_fixed_time,
	:is_finished,
	:is_deleted,
	:sharing_status
)`

/*
	Insert new parcelEvent into DB
*/
func InsertParcelEvents(db *sqlx.DB, events []*ParcelEvent) error {

	var err error

	for _, event := range events {
		_, err = db.NamedExec(parcelEventInsert, event)
		if err != nil {
			return err
		}
		update := `UPDATE parcel_event SET sharing_status = 30 WHERE uid = '` + event.Id + `' AND sharing_status = 10`
		_, err = db.Exec(update)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
	Return SQL with sharing status 20 to the tablet
*/
func GetUnsyncedParcelEventsAsSqlInsert(db *sqlx.DB) (*string, error) {
	rows, err := db.Query("SELECT * FROM parcel_event WHERE sharing_status = 20")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sql := getParcelEventSqlInsert(db, rows)

	return &sql, nil
}

func getParcelEventSqlInsert(db *sqlx.DB, rows *sql.Rows) string {
	var id interface{}
	var createdAt interface{}
	var eventType interface{}
	var parcelUid interface{}
	var ryoseiUid interface{}
	var roomId interface{}
	var name interface{}
	var targetId interface{}
	var note interface{}
	var isAfterPeriodicCheck interface{}
	var isFinished interface{}
	var isDeleted interface{}
	var sharingStatus interface{}

	sql := ""

	for rows.Next() {
		err := rows.Scan(
			&id,
			&createdAt,
			&eventType,
			&parcelUid,
			&ryoseiUid,
			&roomId,
			&name,
			&targetId,
			&note,
			&isAfterPeriodicCheck,
			&isFinished,
			&isDeleted,
			&sharingStatus,
		)

		if err != nil {
			return err.Error()
		}

		query := fmt.Sprintf(
			`INSERT INTO parcel_event(
				uid,
				created_at,
				event_type,
				parcel_uid,
				ryosei_uid,
				room_name,
				ryosei_name,
				target_event_uid,
				note,
				is_after_fixed_time,
				is_finished,
				is_deleted,
				sharing_status
			) VALUES (
				'%s','%s',%d,%v,%v,%v,%v,%v,%v,%d,%d,%d,%d
		);`,
			id,
			createdAt.(time.Time).Format("2006-01-02 15:04:05"),
			eventType,
			nullStringToJsonFormat(parcelUid),
			nullStringToJsonFormat(ryoseiUid),
			nullStringToJsonFormat(roomId),
			nullStringToJsonFormat(name),
			nullStringToJsonFormat(targetId),
			nullStringToJsonFormat(note),
			boolToInt(isAfterPeriodicCheck),
			boolToInt(isFinished),
			boolToInt(isDeleted),
			sharingStatus,
		)
		sql += query
	}
	return sql
}
