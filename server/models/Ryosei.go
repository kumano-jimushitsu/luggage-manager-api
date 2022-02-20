package models

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func (ryosei Ryosei) GetName() string {
	return "ryosei"
}

func (ryosei Ryosei) Uid() string {
	return ryosei.Id
}

func (ryosei *Ryosei) UnmarshalJSON(data []byte) error {

	var record *map[string]interface{}
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	err = setRyosei(ryosei, record)
	if err != nil {
		return err
	}

	return nil
}

func setRyosei(ryosei *Ryosei, record *map[string]interface{}) error {
	ryosei.Id = (*record)["uid"].(string)
	ryosei.RoomID = (*record)["room_name"].(string)
	ryosei.Name = (*record)["ryosei_name"].(string)
	ryosei.Kana = toNullString((*record)["ryosei_name_kana"])
	ryosei.Romaji = toNullString((*record)["ryosei_name_alphabet"])
	ryosei.BlockID = floatToInt((*record)["block_id"].(float64))
	ryosei.SlackID = toNullString((*record)["slack_id"])
	ryosei.Status = floatToInt((*record)["status"].(float64))
	ryosei.CurrentCount = floatToInt((*record)["parcels_current_count"].(float64))
	ryosei.TotalCount = floatToInt((*record)["parcels_total_count"].(float64))
	ryosei.TotalWaitTime = (*record)["parcels_total_waittime"].(string)
	ryosei.LastEventID = toNullString((*record)["last_event_id"])
	ryosei.LastEventDatetime = toNullString((*record)["last_event_datetime"])
	ryosei.CreatedAt = (*record)["created_at"].(string)
	ryosei.UpdatedAt = toNullString((*record)["updated_at"])
	ryosei.SharingStatus = floatToInt((*record)["sharing_status"].(float64))
	ryosei.SharingTime = toNullString((*record)["sharing_time"])
	return nil
}

type Ryosei struct {
	Id                string         `json:"uid" db:"uid"`
	RoomID            string         `json:"room_name" db:"room_name"`
	Name              string         `json:"ryosei_name" db:"ryosei_name"`
	Kana              sql.NullString `json:"ryosei_name_kana" db:"ryosei_name_kana"`
	Romaji            sql.NullString `json:"ryosei_name_alphabet" db:"ryosei_name_alphabet"`
	BlockID           int            `json:"block_id" db:"block_id"`
	SlackID           sql.NullString `json:"slack_id" db:"slack_id"`
	Status            int            `json:"status" db:"status"`
	CurrentCount      int            `json:"parcels_current_count" db:"parcels_current_count"`
	TotalCount        int            `json:"parcels_total_count" db:"parcels_total_count"`
	TotalWaitTime     string         `json:"parcels_total_waittime" db:"parcels_total_waittime"`
	LastEventID       sql.NullString `json:"last_event_id" db:"last_event_id"`
	LastEventDatetime sql.NullString `json:"last_event_datetime" db:"last_event_datetime"`
	CreatedAt         string         `json:"created_at" db:"created_at"`
	UpdatedAt         sql.NullString `json:"updated_at" db:"updated_at"`
	SharingStatus     int            `json:"sharing_status" db:"sharing_status"`
	SharingTime       sql.NullString `json:"sharing_time" db:"sharing_time"`
}

func ParseJsonToRyoseis(raw_json string) ([]*Ryosei, error) {
	var ryoseis []*Ryosei
	err := json.Unmarshal([]byte(raw_json), &ryoseis)
	if err != nil {
		return nil, err
	}
	return ryoseis, err
}

func GetAllRyoseis(db *sqlx.DB) ([]*Ryosei, error) {
	rows, err := db.Query("SELECT * FROM ryosei")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ryoseis, err := getRyoseisFromSqlRows(db, rows)
	if err != nil {
		return nil, err
	}

	return ryoseis, nil
}

func getRyoseisFromSqlRows(db *sqlx.DB, rows *sql.Rows) ([]*Ryosei, error) {

	ryoseis := make([]*Ryosei, 0)
	for rows.Next() {
		ryosei := new(Ryosei)
		err := rows.Scan(
			&ryosei.Id,
			&ryosei.RoomID,
			&ryosei.Name,
			&ryosei.Kana,
			&ryosei.Romaji,
			&ryosei.BlockID,
			&ryosei.SlackID,
			&ryosei.Status,
			&ryosei.CurrentCount,
			&ryosei.TotalCount,
			&ryosei.TotalWaitTime,
			&ryosei.LastEventID,
			&ryosei.LastEventDatetime,
			&ryosei.CreatedAt,
			&ryosei.UpdatedAt,
			&ryosei.SharingStatus,
			&ryosei.SharingTime,
		)
		if err != nil {
			return nil, err
		}
		ryoseis = append(ryoseis, ryosei)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ryoseis, nil
}

var ryoseiInsert string = `
merge into ryosei as old
using
(select
	:uid as uid,
	:room_name as room_name,
	:ryosei_name as ryosei_name,
	:ryosei_name_kana as ryosei_name_kana,
	:ryosei_name_alphabet as ryosei_name_alphabet,
	:block_id as block_id,
	:slack_id as slack_id,
	:status as status,
	:parcels_current_count as parcels_current_count,
	:parcels_total_count as parcels_total_count,
	:parcels_total_waittime as parcels_total_waittime,
	:last_event_id as last_event_id,
	:last_event_datetime as last_event_datetime,
	:created_at as created_at,
	:updated_at as updated_at,
	:sharing_status as sharing_status,
	:sharing_time as sharing_time
) as new
on(
 old.uid=new.uid
)
when matched then
	update set
		uid=new.uid,
		room_name=new.room_name,
		ryosei_name=new.ryosei_name,
		ryosei_name_kana=new.ryosei_name_kana,
		ryosei_name_alphabet=new.ryosei_name_alphabet,
		block_id=new.block_id,
		slack_id=new.slack_id,
		status=new.status,
		parcels_current_count=new.parcels_current_count,
		parcels_total_count=new.parcels_total_count,
		parcels_total_waittime=new.parcels_total_waittime,
		last_event_id=new.last_event_id,
		last_event_datetime=new.last_event_datetime,
		created_at=new.created_at,
		updated_at=new.updated_at,
		sharing_status = 30,
	    sharing_time = getdate()
when not matched then
 insert(
	uid,
	room_name,
	ryosei_name,
	ryosei_name_kana,
	ryosei_name_alphabet,
	block_id,
	slack_id,
	status,
	parcels_current_count,
	parcels_total_count,
	parcels_total_waittime,
	last_event_id,
	last_event_datetime,
	created_at,
	updated_at,
	sharing_status,
	sharing_time
)
 values(
	new.uid,
	new.room_name,
	new.ryosei_name,
	new.ryosei_name_kana,
	new.ryosei_name_alphabet,
	new.block_id,
	new.slack_id,
	new.status,
	new.parcels_current_count,
	new.parcels_total_count,
	new.parcels_total_waittime,
	new.last_event_id,
	new.last_event_datetime,
	new.created_at,
	new.updated_at,
	30,
	getdate()
 );
 `

/*
var ryoseiInsert string = `
INSERT INTO ryosei(
	uid,
	room_name,
	ryosei_name,
	ryosei_name_kana,
	ryosei_name_alphabet,
	block_id,
	slack_id,
	status,
	parcels_current_count,
	parcels_total_count,
	parcels_total_waittime,
	last_event_id,
	last_event_datetime,
	created_at,
	updated_at,
	sharing_status
) VALUES (
	:uid,
	:room_name,
	:ryosei_name,
	:ryosei_name_kana,
	:ryosei_name_alphabet,
	:block_id,
	:slack_id,
	:status,
	:parcels_current_count,
	:parcels_total_count,
	:parcels_total_waittime,
	:last_event_id,
	:last_event_datetime,
	:created_at,
	:updated_at,
	:sharing_status
)`
*/
func InsertRyoseis(db *sqlx.DB, ryoseis []*Ryosei) error {

	var err error

	for _, ryosei := range ryoseis {
		_, err = db.NamedExec(ryoseiInsert, ryosei)
		if err != nil {
			return err
		}
		//↓必要なさそうな気がする
		update := `UPDATE ryosei SET sharing_status = 30 WHERE uid = '` + ryosei.Id + `' AND sharing_status = 10`
		_, err = db.Exec(update)
		// ↑
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func UpdateRyoseis(db *sqlx.DB, ryoseis []*Ryosei) error {
	for _, ryosei := range ryoseis {

		var sql string
		//cratedAt, _ := sqlNullStringToMssqlDateTime(ryosei.CreatedAt)
		updatedAt, _ := sqlNullStringToMssqlDateTime(ryosei.UpdatedAt)

		// sharing_status=11だけどPCにデータが無い時の処理が必要
		count, err := getParcelCountByUid(db, ryosei.Id)
		if err != nil {
			return err
		}

		if count == 0 {
			_, err = db.NamedExec(ryoseiInsert, ryosei)
			if err != nil {
				return err
			}
		} else {
			sql = fmt.Sprintf(`
			UPDATE ryosei
			SET
				uid = '%s',
				room_name = '%s',
				ryosei_name = '%s',
				ryosei_name_kana = '%s',
				ryosei_name_alphabet = '%s',
				block_id = %d,
				slack_id = '%s',
				status = %d,
				parcels_current_count = %d,
				parcels_total_count = %d,
				parcels_total_waittime='%s'
				last_event_id = '%s',
				last_event_datetime = %v,
				created_at = '%s',
				updated_at = %v,
				sharing_status = %d,
			where
				uid = '%s'
			`,
				ryosei.Id,
				ryosei.RoomID,
				ryosei.Name,
				sqlNullStringToJsonFormat(ryosei.Kana),
				sqlNullStringToJsonFormat(ryosei.Romaji),
				ryosei.BlockID,
				sqlNullStringToJsonFormat(ryosei.SlackID),
				ryosei.Status,
				ryosei.CurrentCount,
				ryosei.TotalCount,
				ryosei.TotalWaitTime,
				sqlNullStringToJsonFormat(ryosei.LastEventID),
				sqlNullStringToJsonFormat(ryosei.LastEventDatetime),
				ryosei.CreatedAt,
				updatedAt,
				ryosei.SharingStatus,
				ryosei.Id,
			)
			_, err = db.Exec(sql)
			if err != nil {
				return err
			}
		}

		sql = `UPDATE parcels SET sharing_status = 30 WHERE uid = '` + ryosei.Id + `' AND sharing_status = 11`
		_, err = db.Exec(sql)

		if err != nil {
			return err
		}
	}

	return nil
}
*/
func GetUnsyncedRyoseisAsSqlInsert(db *sqlx.DB) (*string, error) {
	var selectsql string
	selectsql = `SELECT TOP (5) 
	uid,
	room_name,
	ryosei_name,
	ryosei_name_kana,
	ryosei_name_alphabet,
	block_id,
	slack_id,
	status,
	parcels_current_count,
	parcels_total_count,
	parcels_total_waittime,
	last_event_id,
	format(last_event_datetime,'yyyy-MM-dd HH:mm:ss') as last_event_datetime,
	format(created_at,'yyyy-MM-dd HH:mm:ss') as created_at,
	format(updated_at,'yyyy-MM-dd HH:mm:ss') as updated_at,
	sharing_status,
	FORMAT(getdate(),'yyyy/MM/dd HH:mm:ss') as sharing_time
	FROM ryosei WHERE sharing_status = 20`
	//取得時にフォーマットしてしまっているが、go側で以下のようにフォーマットすることもできる
	//(ParcelEvent.goを参照)
	//createdAt.(time.Time).Format("2006-01-02 15:04:05"),
	//直すのが面倒だし、何かの時に使えるかもしれないので残しておく
	rows, err := db.Query(selectsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sql := getSqlRyoseiInsert(db, rows)

	return &sql, nil
}

func getSqlRyoseiInsert(db *sqlx.DB, rows *sql.Rows) string {
	var Id interface{}
	var RoomID interface{}
	var Name interface{}
	var Kana interface{}
	var Romaji interface{}
	var BlockID interface{}
	var SlackID interface{}
	var Status interface{}
	var CurrentCount interface{}
	var TotalCount interface{}
	var TotalWaitTime interface{}
	var LastEventID interface{}
	var LastEventDatetime interface{}
	var CreatedAt interface{}
	var UpdatedAt interface{}
	var SharingStatus interface{}
	var SharingTime interface{}

	sql := ""
	for rows.Next() {
		err := rows.Scan(
			&Id,
			&RoomID,
			&Name,
			&Kana,
			&Romaji,
			&BlockID,
			&SlackID,
			&Status,
			&CurrentCount,
			&TotalCount,
			&TotalWaitTime,
			&LastEventID,
			&LastEventDatetime,
			&CreatedAt,
			&UpdatedAt,
			&SharingStatus,
			&SharingTime,
		)
		if err != nil {
			return err.Error()
		}
		/*
			query := fmt.Sprintf(
				`REPLACE ryosei
				SET
					uid='%s',
					room_name='%s',
					ryosei_name='%s',
					ryosei_name_kana=%v,
					ryosei_name_alphabet=%v,
					block_id=%d,
					slack_id=%v,
					status=%v,
					parcel_current_count=%d,
					parcels_total_count=%d,
					parcels_total_waittime='%s',
					last_event_id=%v,
					last_event_datetime=%v,
					created_at='%s',
					updated_at=%v,
					sharing_status=%v;`,
				Id,
				RoomID,
				Name,
				nullStringToJsonFormat(Kana),
				nullStringToJsonFormat(Romaji),
				BlockID,
				nullStringToJsonFormat(SlackID),
				Status,
				CurrentCount,
				TotalCount,
				TotalWaitTime,
				nullStringToJsonFormat(LastEventID),
				nullStringToJsonFormat(LastEventDatetime),
				CreatedAt,
				nullStringToJsonFormat(UpdatedAt),
				SharingStatus,
			)
		*/
		query := fmt.Sprintf(
			`REPLACE INTO ryosei(
					uid,
					room_name,
					ryosei_name,
					ryosei_name_kana,
					ryosei_name_alphabet,
					block_id,
					slack_id,
					status,
					parcels_current_count,
					parcels_total_count,
					parcels_total_waittime,
					last_event_id,
					last_event_datetime,
					created_at,
					updated_at,
					sharing_status,
					sharing_time
				)VALUES(
					'%s','%s','%s',%v,%v,%d,%v,%d,%d,%d,'%s',%v,%v,'%s',%v,%d,%v
				);`,
			Id,
			RoomID,
			Name,
			nullStringToJsonFormat(Kana),
			nullStringToJsonFormat(Romaji),
			BlockID,
			nullStringToJsonFormat(SlackID),
			Status,
			CurrentCount,
			TotalCount,
			TotalWaitTime,
			nullStringToJsonFormat(LastEventID),
			nullStringToJsonFormat(LastEventDatetime),
			CreatedAt,
			nullStringToJsonFormat(UpdatedAt),
			SharingStatus,
			nullStringToJsonFormat(SharingTime),
		)
		sql += query
	}
	return sql
}

/*
func GetUnsyncedRyoseisAsSqlUpdate(db *sqlx.DB) (*string, error) {
	rows, err := db.Query("SELECT * FROM ryosei WHERE sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sql := getSqlRyoseiUpdate(db, rows)

	return &sql, nil
}

func getSqlRyoseiUpdate(db *sqlx.DB, rows *sql.Rows) string {
	var Id interface{}
	var RoomID interface{}
	var Name interface{}
	var Kana interface{}
	var Romaji interface{}
	var BlockID interface{}
	var SlackID interface{}
	var Status interface{}
	var CurrentCount interface{}
	var TotalCount interface{}
	var TotalWaitTime interface{}
	var LastEventID interface{}
	var LastEventDatetime interface{}
	var CreatedAt interface{}
	var UpdatedAt interface{}
	var SharingStatus interface{}

	sql := ""
	for rows.Next() {
		err := rows.Scan(
			&Id,
			&RoomID,
			&Name,
			&Kana,
			&Romaji,
			&BlockID,
			&SlackID,
			&Status,
			&CurrentCount,
			&TotalCount,
			&TotalWaitTime,
			&LastEventID,
			&LastEventDatetime,
			&CreatedAt,
			&UpdatedAt,
			&SharingStatus,
		)
		if err != nil {
			return err.Error()
		}
		query := fmt.Sprintf(`
			UPDATE ryosei
				SET
					uid = '%s',
					room_name = '%s',
					ryosei_name = '%s',
					ryosei_name_kana = '%s',
					ryosei_name_alphabet = '%s',
					block_id=%d,
					slack_id = '%s',
					status=%d,
					parcel_current_count=%d,
					parcels_total_count=%d,
					parcels_total_waittime=%d,
					last_event_id = '%s',
					last_event_datetime = '%s'
					created_at = '%s',
					updated_at = %v,
					sharing_status=%d
				WHERE
					uid='%s'
			;`,
			Id,
			RoomID,
			Name,
			nullStringToJsonFormat(Kana),
			nullStringToJsonFormat(Romaji),
			BlockID,
			nullStringToJsonFormat(SlackID),
			Status,
			CurrentCount,
			TotalCount,
			TotalWaitTime,
			nullStringToJsonFormat(LastEventID),
			nullTimeToJsonFormat(LastEventDatetime),
			CreatedAt.(time.Time).Format("2006-01-02 15:04:05"),
			UpdatedAt,
			SharingStatus,
			Id,
		)
		sql += query
	}
	return sql
}
*/
// [Deprecated] Make a long sql insert from data in DB
func GetRyoseiSeedingSql(db *sqlx.DB) (string, error) {
	rows, err := db.Query("SELECT * FROM ryosei")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var sql string

	var id interface{}
	var roomID interface{}
	var name interface{}
	var kana interface{}
	var romaji interface{}
	var blockID interface{}
	var slackID interface{}
	var status interface{}
	var currentCount interface{}
	var totalCount interface{}
	var totalWaitTime interface{}
	var lastEventID interface{}
	var lastEventDatetime interface{}
	var createdAt interface{}
	var updateAt interface{}
	var sharingStatus interface{}

	for rows.Next() {
		err := rows.Scan(
			&id,
			&roomID,
			&name,
			&kana,
			&romaji,
			&blockID,
			&slackID,
			&status,
			&currentCount,
			&totalCount,
			&totalWaitTime,
			&lastEventID,
			&lastEventDatetime,
			&createdAt,
			&updateAt,
			&sharingStatus,
		)

		if err != nil {
			return "", err
		}

		query := fmt.Sprintf(
			`INSERT INTO ryosei(
				uid,
				room_name,
				ryosei_name,
				ryosei_name_kana,
				ryosei_name_alphabet,
				block_id,
				slack_id,
				status,
				parcels_current_count,
				parcels_total_count,
				parcels_total_waittime,
				last_event_id,
				last_event_datetime,
				created_at,
				updated_at,
				sharing_status
			) VALUES(
				%s,%s,%s,%s,%s,%d,%s,%d,%d,%d,%s,%s,%s,%s,%s,%d
		);`,
			id,
			roomID,
			name,
			nullStringToJsonFormat(kana),
			nullStringToJsonFormat(romaji),
			blockID,
			nullStringToJsonFormat(slackID),
			status,
			currentCount,
			totalCount,
			totalWaitTime,
			nullStringToJsonFormat(lastEventID),
			nullStringToJsonFormat(lastEventDatetime),
			createdAt,
			nullStringToJsonFormat(updateAt),
			sharingStatus,
		)
		sql += query
	}
	return sql, nil
}

/*
	[Depreated] Seed database from csv ryosei data
*/
func GetRyoseiSeedingCsv(db *sqlx.DB) {
	rows, err := db.Query("SELECT * FROM ryosei")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id interface{}
	var roomID interface{}
	var name interface{}
	var kana interface{}
	var romaji interface{}
	var blockID interface{}
	var slackID interface{}
	var status interface{}
	var currentCount interface{}
	var totalCount interface{}
	var totalWaitTime interface{}
	var lastEventID interface{}
	var lastEventDatetime interface{}
	var createdAt interface{}
	var updateAt interface{}
	var sharingStatus interface{}

	file, err := os.Create("寮生.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cw := csv.NewWriter(file)
	defer cw.Flush()

	for rows.Next() {
		err := rows.Scan(
			&id,
			&roomID,
			&name,
			&kana,
			&romaji,
			&blockID,
			&slackID,
			&status,
			&currentCount,
			&totalCount,
			&totalWaitTime,
			&lastEventID,
			&lastEventDatetime,
			&createdAt,
			&updateAt,
			&sharingStatus,
		)

		if err != nil {
			log.Fatal(err)
		}

		createdAt, ok := createdAt.(time.Time)
		if ok == false {
			panic("Type assertion of createdAt into time.Time failed")
		}

		col := []string{
			id.(string),
			roomID.(string),
			name.(string),
			kana.(string),
			romaji.(string),
			fmt.Sprint(blockID.(int64)),
			nullStringToJsonFormat(slackID),
			fmt.Sprint(status.(int64)),
			fmt.Sprint(currentCount.(int64)),
			fmt.Sprint(totalCount.(int64)),
			totalWaitTime.(string),
			nullStringToJsonFormat(lastEventID),
			nullTimeToJsonFormat(lastEventDatetime),
			nullTimeToJsonFormat(createdAt),
			nullTimeToJsonFormat(updateAt),
			fmt.Sprint(sharingStatus.(int64)),
		}

		cw.Write(col)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

func IncrementParcelCount(db *sqlx.DB, parcel Parcel) error {
	ownerId := parcel.OwnerID
	sql := fmt.Sprintf("UPDATE ryosei SET parcels_current_count = parcels_current_count + 1, parcels_total_count = parcels_total_count + 1 WHERE uid = '%s'", ownerId)
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func IncrementParcelCountSql(db *sqlx.DB, ownerId string) string {
	sql := fmt.Sprintf("UPDATE ryosei SET parcels_current_count = parcels_current_count + 1, parcels_total_count = parcels_total_count + 1 WHERE ryosei_name = '%s';", ownerId)
	return sql
}
