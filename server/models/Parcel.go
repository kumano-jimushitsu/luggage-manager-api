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
func (parcel Parcel) GetName() string {
	return "parcels"
}

/*
	Implement Unmarshaler interface
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
	parcel.OwnerRyoseiName = (*record)["owner_ryosei_name"].(string)
	parcel.RegisteredAt = (*record)["register_datetime"].(string)
	parcel.RegisteredStaffID = (*record)["register_staff_uid"].(string)
	parcel.RegisteredStaffRoomName = (*record)["register_staff_room_name"].(string)
	parcel.RegisteredStaffName = (*record)["register_staff_ryosei_name"].(string)
	parcel.Placement = floatToInt((*record)["placement"].(float64))
	parcel.Fragile = floatToInt((*record)["fragile"].(float64))
	parcel.IsReleased = floatToInt((*record)["is_released"].(float64))
	parcel.ReleasedAgentID = toNullString((*record)["release_agent_uid"])
	parcel.ReleasedAt = toNullString((*record)["release_datetime"])
	parcel.ReleasedStaffID = toNullString((*record)["released_staff_uid"])
	parcel.ReleasedStaffRoomID = toNullString((*record)["released_staff_room_name"])
	parcel.ReleasedStaffName = toNullString((*record)["released_staff_ryosei_name"])
	parcel.CheckedCount = floatToInt((*record)["checked_count"].(float64))
	parcel.IsLost = floatToInt((*record)["is_lost"].(float64))
	parcel.LostAt = toNullString((*record)["lost_datetime"])
	parcel.IsReturned = floatToInt((*record)["is_returned"].(float64))
	parcel.ReturnedAt = toNullString((*record)["returned_datetime"])
	parcel.IsOperationError = floatToInt((*record)["is_operation_error"].(float64))
	parcel.OperationErrorType = toNullInt32((*record)["operation_error_type"])
	parcel.Description = toNullString((*record)["note"])
	parcel.IsDeleted = floatToInt((*record)["is_deleted"].(float64))
	parcel.SharingStatus = floatToInt((*record)["sharing_status"].(float64))
	//parcel.SharingTime = toNullString((*record)["sharing_status"])
	return nil
}

type Parcel struct {
	Id                      string         `json:"uid" db:"uid"`
	OwnerID                 string         `json:"owner_uid" db:"owner_uid"`
	OwnerRoomID             string         `json:"owner_room_name" db:"owner_room_name"`
	OwnerRyoseiName         string         `json:"owner_ryosei_name" db:"owner_ryosei_name"`
	RegisteredAt            string         `json:"register_datetime" db:"register_datetime"`
	RegisteredStaffID       string         `json:"register_staff_uid" db:"register_staff_uid"`
	RegisteredStaffRoomName string         `json:"register_staff_room_name" db:"register_staff_room_name"`
	RegisteredStaffName     string         `json:"register_staff_ryosei_name" db:"register_staff_ryosei_name"`
	Placement               int            `json:"placement" db:"placement"`
	Fragile                 int            `json:"fragile" db:"fragile"`
	IsReleased              int            `json:"is_released" db:"is_released"`
	ReleasedAgentID         sql.NullString `json:"release_agent_uid" db:"release_agent_uid"`
	ReleasedAt              sql.NullString `json:"release_datetime" db:"release_datetime"`
	ReleasedStaffID         sql.NullString `json:"release_staff_uid" db:"release_staff_uid"`
	ReleasedStaffRoomID     sql.NullString `json:"release_staff_room_name" db:"release_staff_room_name"`
	ReleasedStaffName       sql.NullString `json:"release_staff_ryosei_name" db:"release_staff_ryosei_name"`
	CheckedCount            int            `json:"checked_count" db:"checked_count"`
	IsLost                  int            `json:"is_lost" db:"is_lost"`
	LostAt                  sql.NullString `json:"lost_datetime" db:"lost_datetime"`
	IsReturned              int            `json:"is_returned" db:"is_returned"`
	ReturnedAt              sql.NullString `json:"returned_datetime" db:"returned_datetime"`
	IsOperationError        int            `json:"is_operation_error" db:"is_operation_error"`
	OperationErrorType      sql.NullInt32  `json:"operation_error_type" db:"operation_error_type"`
	Description             sql.NullString `json:"note" db:"note"`
	IsDeleted               int            `json:"is_deleted" db:"is_deleted"`
	SharingStatus           int            `json:"sharing_status" db:"sharing_status"`
	//SharingTime             sql.NullString `json:"sharing_time" db:"sharing_time"`
}

/*
	Parse json and return the slice of struct Parcel
*/
func ParseJsonToParcels(raw_json string) ([]*Parcel, error) {
	var parcels []*Parcel
	err := json.Unmarshal([]byte(raw_json), &parcels)
	if err != nil {
		return nil, err
	}
	return parcels, err
}

/*
	Return all records in the parcels table
*/
func GetAllParcels(db *sqlx.DB) ([]*Parcel, error) {
	rows, err := db.Query("SELECT * FROM parcels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parcels, err := getParcelsFromSqlRows(db, rows)
	if err != nil {
		return nil, err
	}

	return parcels, nil
}

func getParcelsFromSqlRows(db *sqlx.DB, rows *sql.Rows) ([]*Parcel, error) {

	parcels := make([]*Parcel, 0)
	for rows.Next() {
		var isReleased bool
		var isLost bool
		var isReturned bool
		var isOperationError bool
		var isDeleted bool
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
			&isReleased,
			&parcel.ReleasedAgentID,
			&parcel.ReleasedAt,
			&parcel.ReleasedStaffID,
			&parcel.ReleasedStaffRoomID,
			&parcel.ReleasedStaffName,
			&parcel.CheckedCount,
			&isLost,
			&parcel.LostAt,
			&isReturned,
			&parcel.ReturnedAt,
			&isOperationError,
			&parcel.OperationErrorType,
			&parcel.Description,
			&isDeleted,
			&parcel.SharingStatus,
		)
		if err != nil {
			return nil, err
		}
		parcel.IsReleased = boolToInt(isReleased)
		parcel.IsLost = boolToInt(isLost)
		parcel.IsReturned = boolToInt(isReturned)
		parcel.IsOperationError = boolToInt(isOperationError)
		parcel.IsDeleted = boolToInt(isDeleted)

		parcels = append(parcels, parcel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parcels, nil
}

var parcelInsert string = `
merge into parcels as old
using
(select
	:uid as uid,
	:owner_uid as owner_uid,
	:owner_room_name as owner_room_name,
	:owner_ryosei_name as owner_ryosei_name,
	:register_datetime as register_datetime,
	:register_staff_uid as register_staff_uid,
	:register_staff_room_name as register_staff_room_name,
	:register_staff_ryosei_name as register_staff_ryosei_name,
	:placement as placement,
	:fragile as fragile,
	:is_released as is_released,
	:release_agent_uid as release_agent_uid,
	:release_datetime as release_datetime,
	:release_staff_uid as release_staff_uid,
	:release_staff_room_name as release_staff_room_name,
	:release_staff_ryosei_name as release_staff_ryosei_name,
	:checked_count as checked_count,
	:is_lost as is_lost,
	:lost_datetime as lost_datetime,
	:is_returned as is_returned,
	:returned_datetime as returned_datetime,
	:is_operation_error as is_operation_error,
	:operation_error_type as operation_error_type,
	:note as note,
	:is_deleted as is_deleted,
	:sharing_status as sharing_status
) as new
on(
	old.uid=new.uid
)
when matched then
update set
 	uid = new.uid,
 	owner_uid = new.owner_uid,
 	owner_room_name = new.owner_room_name,
 	owner_ryosei_name = new.owner_ryosei_name,
 	register_datetime = new.register_datetime,
 	register_staff_uid = new.register_staff_uid,
 	register_staff_room_name = new.register_staff_room_name,
 	register_staff_ryosei_name = new.register_staff_ryosei_name,
 	placement = new.placement,
 	fragile = new.fragile,
 	is_released = new.is_released,
 	release_agent_uid = new.release_agent_uid,
 	release_datetime = new.release_datetime,
 	release_staff_uid = new.release_staff_uid,
 	release_staff_room_name = new.release_staff_room_name,
 	release_staff_ryosei_name = new.release_staff_ryosei_name,
 	checked_count = new.checked_count,
 	is_lost = new.is_lost,
 	lost_datetime = new.lost_datetime,
 	is_returned = new.is_returned,
 	returned_datetime = new.returned_datetime,
 	is_operation_error = new.is_operation_error,
 	operation_error_type = new.operation_error_type,
 	note = new.note,
 	is_deleted = new.is_deleted,
 	sharing_status = new.sharing_status
when not matched then
 insert(
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
 )
 values(
	new.uid,
	new.owner_uid,
	new.owner_room_name,
	new.owner_ryosei_name,
	new.register_datetime,
	new.register_staff_uid,
	new.register_staff_room_name,
	new.register_staff_ryosei_name,
	new.placement,
	new.fragile,
	new.is_released,
	new.release_agent_uid,
	new.release_datetime,
	new.release_staff_uid,
	new.release_staff_room_name,
	new.release_staff_ryosei_name,
	new.checked_count,
	new.is_lost,
	new.lost_datetime,
	new.is_returned,
	new.returned_datetime,
	new.is_operation_error,
	new.operation_error_type,
	new.note,
	new.is_deleted,
	new.sharing_status
 );
`

/*
var parcelInsert string = `
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
) VALUES (
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
*/
/*
	Insert new parcels into DB
*/
func InsertParcels(db *sqlx.DB, parcels []*Parcel) error {

	var err error

	for _, parcel := range parcels {
		_, err = db.NamedExec(parcelInsert, parcel)
		if err != nil {
			return err
		}
		update := `UPDATE parcels SET sharing_status = 30 WHERE uid = '` + parcel.Id + `' AND sharing_status = 10`
		_, err = db.Exec(update)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
	Update records in the table with the latest parcels
*/
func UpdateParcels(db *sqlx.DB, parcels []*Parcel) error {
	for _, parcel := range parcels {

		var sql string

		registeredAt, _ := stringToMssqlDateTime(parcel.RegisteredAt)
		releasedAt, _ := sqlNullStringToMssqlDateTime(parcel.ReleasedAt)
		lostAt, _ := sqlNullStringToMssqlDateTime(parcel.LostAt)
		returnedAt, _ := sqlNullStringToMssqlDateTime(parcel.ReturnedAt)

		// sharing_status=11だけどPCにデータが無い時の処理が必要
		count, err := getParcelCountByUid(db, parcel.Id)
		if err != nil {
			return err
		}

		if count == 0 {
			_, err = db.NamedExec(parcelInsert, parcel)
			if err != nil {
				return err
			}
		} else {
			sql = fmt.Sprintf(`
			UPDATE parcels
			SET
				uid = '%s',
				owner_uid = '%s',
				owner_room_name = '%s',
				owner_ryosei_name = '%s',
				register_datetime = '%s',
				register_staff_uid = '%s',
				register_staff_room_name = '%s',
				register_staff_ryosei_name = '%s',
				placement = %d,
				fragile = %d,
				is_released = %d,
				release_agent_uid = %v,
				release_datetime = %v,
				release_staff_uid = %v,
				release_staff_room_name = %v,
				release_staff_ryosei_name = %v,
				checked_count = %d,
				is_lost = %d,
				lost_datetime = %v,
				is_returned = %d,
				returned_datetime = %v,
				is_operation_error = %d,
				operation_error_type = %v,
				note = %v,
				is_deleted = %d,
				sharing_status = %d
			WHERE
				uid = '%s'
			`,
				parcel.Id,
				parcel.OwnerID,
				parcel.OwnerRoomID,
				parcel.OwnerRyoseiName,
				registeredAt,
				parcel.RegisteredStaffID,
				parcel.RegisteredStaffRoomName,
				parcel.RegisteredStaffName,
				parcel.Placement,
				parcel.Fragile,
				boolToInt(parcel.IsReleased),
				sqlNullStringToJsonFormat(parcel.ReleasedAgentID),
				releasedAt,
				sqlNullStringToJsonFormat(parcel.ReleasedStaffID),
				sqlNullStringToJsonFormat(parcel.ReleasedStaffRoomID),
				sqlNullStringToJsonFormat(parcel.ReleasedStaffName),
				parcel.CheckedCount,
				boolToInt(parcel.IsLost),
				lostAt,
				boolToInt(parcel.IsReturned),
				returnedAt,
				boolToInt(parcel.IsOperationError),
				nullInt32ToJsonFormat(parcel.OperationErrorType),
				sqlNullStringToJsonFormat(parcel.Description),
				boolToInt(parcel.IsDeleted),
				parcel.SharingStatus,
				parcel.Id,
			)
			_, err = db.Exec(sql)
			if err != nil {
				return err
			}
		}
		sql = `UPDATE parcels SET sharing_status = 30 WHERE uid = '` + parcel.Id + `' AND sharing_status = 11`
		_, err = db.Exec(sql)

		if err != nil {
			return err
		}
	}

	return nil
}

func getParcelCountByUid(db *sqlx.DB, uid string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT * FROM parcels WHERE uid = '%s'", uid)
	rows, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		count++
	}
	return count, nil
}

/*
	Return SQL with sharing status 20 to the tablet
*/
func GetUnsyncedParcelsAsSqlInsert(db *sqlx.DB) (*string, error) {
	rows, err := db.Query("SELECT TOP(50) * FROM parcels WHERE sharing_status = 20")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sql := getSqlParcelInsert(db, rows)

	return &sql, nil
}

func getSqlParcelInsert(db *sqlx.DB, rows *sql.Rows) string {
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

	sql := ""
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
			`REPLACE INTO parcels(
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
				'%s','%s','%s','%s','%s','%s','%s','%s',%d,%d,
				%d,%v,%v,%v,%v,%v,%d,%d,%v,%d,
				%v,%d,%v,%v,%d,%d
		);`,
			id,
			ownerID,
			ownerRoomID,
			ownerRyoseiName,
			registeredAt.(time.Time).Format("2006-01-02 15:04:05"),
			registeredStaffID,
			registeredStaffRoomName,
			registeredStaffName,
			placement,
			fragile,
			boolToInt(isReleased),
			nullStringToJsonFormat(releasedAgentID),
			nullTimeToJsonFormat(releasedAt),
			nullStringToJsonFormat(releasedStaffID),
			nullStringToJsonFormat(releasedStaffRoomID),
			nullStringToJsonFormat(releasedStaffName),
			checkedCount,
			boolToInt(isLost),
			nullTimeToJsonFormat(lostAt),
			boolToInt(isReturned),
			nullTimeToJsonFormat(returnedAt),
			boolToInt(isOperationError),
			nullInt32ToJsonFormat(operationErrorType),
			nullStringToJsonFormat(description),
			boolToInt(isDeleted),
			sharingStatus,
		)
		sql += query
	}
	return sql
}

/*
	Return SQL with sharing status 20 and 21 to the tablet
*/
func GetUnsyncedParcelsAsSqlUpdate(db *sqlx.DB) (*string, error) {
	rows, err := db.Query("SELECT * FROM parcels WHERE sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sql := getSqlParcelsUpdate(db, rows)

	return &sql, nil
}

func getSqlParcelsUpdate(db *sqlx.DB, rows *sql.Rows) string {
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

	sql := ""
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

		query := fmt.Sprintf(`
			UPDATE parcels
			SET
				uid = '%s',
				owner_uid = '%s',
				owner_room_name = '%s',
				owner_ryosei_name = '%s',
				register_datetime = '%s',
				register_staff_uid = '%s',
				register_staff_room_name = '%s',
				register_staff_ryosei_name = '%s',
				placement = %d,
				fragile = %d,
				is_released = %d,
				release_agent_uid = %v,
				release_datetime = %v,
				release_staff_uid = %v,
				release_staff_room_name = %v,
				release_staff_ryosei_name = %v,
				checked_count = %d,
				is_lost = %d,
				lost_datetime = %v,
				is_returned = %d,
				returned_datetime = %v,
				is_operation_error = %d,
				operation_error_type = %v,
				note = %v,
				is_deleted = %d,
				sharing_status = %d
			WHERE
				uid = '%s'
			;`,
			id,
			ownerID,
			ownerRoomID,
			ownerRyoseiName,
			registeredAt.(time.Time).Format("2006-01-02 15:04:05"),
			registeredStaffID,
			registeredStaffRoomName,
			registeredStaffName,
			placement,
			fragile,
			boolToInt(isReleased),
			nullStringToJsonFormat(releasedAgentID),
			nullTimeToJsonFormat(releasedAt),
			nullStringToJsonFormat(releasedStaffID),
			nullStringToJsonFormat(releasedStaffRoomID),
			nullStringToJsonFormat(releasedStaffName),
			checkedCount,
			boolToInt(isLost),
			nullTimeToJsonFormat(lostAt),
			boolToInt(isReturned),
			nullTimeToJsonFormat(returnedAt),
			boolToInt(isOperationError),
			nullInt32ToJsonFormat(operationErrorType),
			nullStringToJsonFormat(description),
			boolToInt(isDeleted),
			sharingStatus,
			id,
		)
		sql += query
	}
	return sql
}
