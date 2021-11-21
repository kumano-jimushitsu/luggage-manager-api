package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Parcel struct {
	Id                      string         `json:"uid" db:"uid"`
	OwnerID                 string         `json:"owner_uid" db:"owner_uid"`
	OwnerRoomID             string         `json:"owner_room_name" db:"owner_room_name"`
	OwnerRyoseiName         string         `json:"owner_parcels_name" db:"owner_ryosei_name"`
	RegisteredAt            string         `json:"register_datetime" db:"register_datetime"`
	RegisteredStaffID       string         `json:"register_staff_uid" db:"register_staff_uid"`
	RegisteredStaffRoomName string         `json:"register_staff_room_name" db:"register_staff_room_name"`
	RegisteredStaffName     string         `json:"register_staff_parcels_name" db:"register_staff_ryosei_name"`
	Placement               int            `json:"placement" db:"placement"`
	Fragile                 bool           `json:"fragile" db:"fragile"`
	IsReleased              bool           `json:"is_released" db:"is_released"`
	ReleasedAgentID         sql.NullString `json:"release_agent_uid" db:"release_agent_uid"`
	ReleasedAt              sql.NullString `json:"release_datetime" db:"release_datetime"`
	ReleasedStaffID         sql.NullString `json:"release_staff_uid" db:"release_staff_uid"`
	ReleasedStaffRoomID     sql.NullString `json:"release_staff_room_name" db:"release_staff_room_name"`
	ReleasedStaffName       sql.NullString `json:"release_staff_parcels_name" db:"release_staff_ryosei_name"`
	CheckedCount            int            `json:"checked_count" db:"checked_count"`
	IsLost                  bool           `json:"is_lost" db:"is_lost"`
	LostAt                  sql.NullString `json:"lost_datetime" db:"lost_datetime"`
	IsReturned              bool           `json:"is_returned" db:"is_returned"`
	ReturnedAt              sql.NullString `json:"returned_datetime" db:"returned_datetime"`
	IsOperationError        bool           `json:"is_operation_error" db:"is_operation_error"`
	OperationErrorType      sql.NullInt32  `json:"operation_error_type" db:"operation_error_type"`
	Description             sql.NullString `json:"note" db:"note"`
	IsDeleted               bool           `json:"is_deleted" db:"is_deleted"`
	SharingStatus           int            `json:"sharing_status" db:"sharing_status"`
}

/*
	Go will unmarshal a struct if it implements Unmarshaler interface
	https://golang.org/pkg/encoding/json/#Unmarshaler
*/
func (parcel *Parcel) UnmarshalJSON(data []byte) error {

	// get parcels
	var record *map[string]interface{}
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	// parse number(float64) type to bool
	err = setParcel(parcel, record)
	if err != nil {
		return err
	}

	return err
}

func setParcel(parcel *Parcel, record *map[string]interface{}) error {
	parcel.Id = (*record)["uid"].(string)
	parcel.OwnerID = (*record)["owner_uid"].(string)
	parcel.OwnerRoomID = (*record)["owner_room_name"].(string)
	parcel.OwnerRyoseiName = (*record)["owner_parcels_name"].(string)
	parcel.RegisteredAt = (*record)["register_datetime"].(string)
	parcel.RegisteredStaffID = (*record)["register_staff_uid"].(string)
	parcel.RegisteredStaffRoomName = (*record)["register_staff_room_name"].(string)
	parcel.RegisteredStaffName = (*record)["register_staff_parcels_name"].(string)
	parcel.Placement = floatToInt((*record)["placement"].(float64))
	parcel.Fragile = floatToBool((*record)["fragile"].(float64))
	parcel.IsReleased = floatToBool((*record)["is_released"].(float64))
	parcel.ReleasedAgentID = toNullString((*record)["release_agent_uid"])
	parcel.ReleasedAt = toNullString((*record)["release_datetime"])
	parcel.ReleasedStaffID = toNullString((*record)["released_staff_uid"])
	parcel.ReleasedStaffRoomID = toNullString((*record)["released_staff_room_name"])
	parcel.ReleasedStaffName = toNullString((*record)["released_staff_parcels_name"])
	parcel.CheckedCount = floatToInt((*record)["checked_count"].(float64))
	parcel.IsLost = floatToBool((*record)["is_lost"].(float64))
	parcel.LostAt = toNullString((*record)["lost_datetime"])
	parcel.IsReturned = floatToBool((*record)["is_returned"].(float64))
	parcel.ReturnedAt = toNullString((*record)["returned_datetime"])
	parcel.IsOperationError = floatToBool((*record)["is_operation_error"].(float64))
	parcel.OperationErrorType = toNullInt32((*record)["operation_error_type"])
	parcel.Description = toNullString((*record)["note"])
	parcel.IsDeleted = floatToBool((*record)["is_deleted"].(float64))
	parcel.IsDeleted = floatToBool((*record)["is_deleted"].(float64))
	parcel.SharingStatus = floatToInt((*record)["sharing_status"].(float64))

	return nil
}

func floatToBool(val float64) bool {
	if val == float64(0) {
		return false
	} else {
		return true
	}
}

func floatToInt(val float64) int {
	return int(val)
}

func nilToEmptyString(val interface{}) string {
	if val == nil {
		return ""
	} else {
		return val.(string)
	}
}

func nilToInt(val interface{}) int {
	if val == nil {
		return 0
	} else {
		return val.(int)
	}
}

func toNullString(val interface{}) sql.NullString {
	if val == nil {
		return sql.NullString{"", false}
	} else {
		return sql.NullString{val.(string), true}
	}
}

func toNullInt32(val interface{}) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{0, false}
	} else {
		return sql.NullInt32{val.(int32), true}
	}
}

func GetAllParcels(db *sqlx.DB) ([]*Parcel, error) {
	rows, err := db.Query("SELECT * FROM parcels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parcels, err := getParcelsFromSqlRows(db)
	if err != nil {
		return nil, err
	}

	return parcels, nil
}

func GetUnsyncedParcelsAsSqlInsert(db *sqlx.DB) (*string, error) {
	rows, err := db.Query("SELECT * FROM parcels WHERE sharing_status = 20 OR sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	json := getSqlInsert(db, rows)

	return &json, nil
}

func getParcelsFromSqlRows(db *sqlx.DB) ([]*Parcel, error) {
	rows, err := db.Query("SELECT * FROM parcels WHERE sharing_status = 20 OR sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parcels := make([]*Parcel, 0)
	for rows.Next() {
		parcel := new(Parcel)
		err := rows.Scan(
			&parcel.Id,
			&parcel.OwnerID,
			&parcel.OwnerRoomID,
			&parcel.OwnerRyoseiName,
			&parcel.RegisteredAt,
			&parcel.RegisteredStaffID,
			&parcel.RegisteredStaffRoomName,
			&parcel.RegisteredStaffName,
			&parcel.Placement,
			&parcel.Fragile,
			&parcel.IsReleased,
			&parcel.ReleasedAgentID,
			&parcel.ReleasedAt,
			&parcel.ReleasedStaffID,
			&parcel.ReleasedStaffRoomID,
			&parcel.ReleasedStaffName,
			&parcel.CheckedCount,
			&parcel.IsLost,
			&parcel.LostAt,
			&parcel.IsReturned,
			&parcel.ReturnedAt,
			&parcel.IsOperationError,
			&parcel.OperationErrorType,
			&parcel.Description,
			&parcel.IsDeleted,
			&parcel.SharingStatus,
		)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, parcel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parcels, nil
}

func getSqlInsert(db *sqlx.DB, rows *sql.Rows) string {
	var id interface{}
	var ownerID interface{}
	var ownerRoomID interface{}
	var ownerRyoseiName interface{}
	var registeredAt interface{}
	var registeredStaffID interface{}
	var registeredStaffRoomName interface{}
	var registeredStaffName interface{}
	var placement interface{}
	var fragile interface{}
	var isReleased interface{}
	var releasedAgentID interface{}
	var releasedAt interface{}
	var releasedStaffID interface{}
	var releasedStaffRoomID interface{}
	var releasedStaffName interface{}
	var checkedCount interface{}
	var isLost interface{}
	var lostAt interface{}
	var isReturned interface{}
	var returnedAt interface{}
	var isOperationError interface{}
	var operationErrorType interface{}
	var description interface{}
	var isDeleted interface{}
	var sharingStatus interface{}

	json := ""
	for rows.Next() {
		err := rows.Scan(
			&id,
			&ownerID,
			&ownerRoomID,
			&ownerRyoseiName,
			&registeredAt,
			&registeredStaffID,
			&registeredStaffRoomName,
			&registeredStaffName,
			&placement,
			&fragile,
			&isReleased,
			&releasedAgentID,
			&releasedAt,
			&releasedStaffID,
			&releasedStaffRoomID,
			&releasedStaffName,
			&checkedCount,
			&isLost,
			&lostAt,
			&isReturned,
			&returnedAt,
			&isOperationError,
			&operationErrorType,
			&description,
			&isDeleted,
			&sharingStatus,
		)
		if err != nil {
			return err.Error()
		}
		query := fmt.Sprintf(
			`INSERT INTO parcels(
				uid,
				owner_uid,
				owner_room_name,
				owner_ryosei_name,
				register_datetime,
				register_staff_uid,
				register_staff_room_name,
				register_staff_ryosei_name,
				placement,
				fragile,
				is_released,
				release_agent_uid,
				release_datetime,
				release_staff_uid,
				release_staff_room_name,
				release_staff_ryosei_name,
				checked_count,
				is_lost,
				lost_datetime,
				is_returned,
				returned_datetime,
				is_operation_error,
				operation_error_type,
				note,
				is_deleted,
				sharing_status
			) VALUES(
				"%v","%v","%v","%v","%v","%v","%v","%v","%v","%v",
				"%v","%v","%v","%v","%v","%v","%v","%v","%v","%v",
				"%v","%v","%v","%v","%v","%v"
		);`,
			id,
			ownerID,
			ownerRoomID,
			ownerRyoseiName,
			registeredAt,
			registeredStaffID,
			registeredStaffRoomName,
			registeredStaffName,
			placement,
			fragile,
			isReleased,
			releasedAgentID,
			releasedAt,
			releasedStaffID,
			releasedStaffRoomID,
			releasedStaffName,
			checkedCount,
			isLost,
			lostAt,
			isReturned,
			returnedAt,
			isOperationError,
			operationErrorType,
			description,
			isDeleted,
			sharingStatus,
		)
		json += query
	}
	return json
}

func InsertParcels(db *sqlx.DB, parcels []*Parcel) error {

	var err error

	query := `
	INSERT INTO parcels(
		uid,
		owner_uid,
		owner_room_name,
		owner_ryosei_name,
		register_datetime,
		register_staff_uid,
		register_staff_room_name,
		register_staff_ryosei_name,
		placement,
		fragile,
		is_released,
		release_agent_uid,
		release_datetime,
		release_staff_uid,
		release_staff_room_name,
		release_staff_ryosei_name,
		checked_count,
		is_lost,
		lost_datetime,
		is_returned,
		returned_datetime,
		is_operation_error,
		operation_error_type,
		note,
		is_deleted,
		sharing_status
	) VALUES(
		:uid,
		:owner_uid,
		:owner_room_name,
		:owner_ryosei_name,
		:register_datetime,
		:register_staff_uid,
		:register_staff_room_name,
		:register_staff_ryosei_name,
		:placement,
		:fragile,
		:is_released,
		:release_agent_uid,
		:release_datetime,
		:release_staff_uid,
		:release_staff_room_name,
		:release_staff_ryosei_name,
		:checked_count,
		:is_lost,
		:lost_datetime,
		:is_returned,
		:returned_datetime,
		:is_operation_error,
		:operation_error_type,
		:note,
		:is_deleted,
		:sharing_status
	)`

	for _, parcel := range parcels {
		_, err = db.NamedExec(query, parcel)
		if err != nil {
			return err
		}
	}

	return nil
}
