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

func (ryosei *Ryosei) UnmarshalJSON(data []byte) error {

	var record *map[string]interface{}
	err := json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	return nil
}

type Ryosei struct {
	Id                string         `json:"uid" db:"uid"`
	RoomID            string         `json:"room_name" db:"room_name"`
	Name              string         `json:"ryosei_name" db:"ryosei_name"`
	Kana              string         `json:"ryosei_name_kana" db:"ryosei_name_kana"`
	Romaji            string         `json:"ryosei_name_alphabet" db:"ryosei_name_alphabet"`
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
		)
		if err != nil {
			return nil, err
		}
		ryoseis = append(ryoseis, ryosei)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ryoseis, nil
}

func getRyoseiFromSqlRows(db *sqlx.DB) ([]*Ryosei, error) {
	rows, err := db.Query("SELECT * FROM ryosei WHERE sharing_status = 20 OR sharing_status = 21")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func setRyosei(ryosei *Ryosei, record *map[string]interface{}) error {
	ryosei.Id = (*record)["uid"].(string)
	ryosei.RoomID = (*record)["room_name"].(string)
	ryosei.Name = (*record)["ryosei_name"].(string)
	ryosei.Kana = (*record)["ryosei_name_kana"].(string)
	ryosei.Romaji = (*record)["ryosei_alphabet"].(string)
	ryosei.BlockID = floatToInt((*record)["block_id"].(float64))
	ryosei.SlackID = toNullString((*record)["slack_id"].(string))
	ryosei.Status = floatToInt((*record)["status"].(float64))
	ryosei.CurrentCount = floatToInt((*record)["parcels_current_count"].(float64))
	ryosei.TotalCount = floatToInt((*record)["parcels_total_count"].(float64))
	ryosei.TotalWaitTime = (*record)["parcels_total_waittime"].(string)
	ryosei.LastEventID = toNullString((*record)["last_event_id"].(string))
	ryosei.LastEventDatetime = toNullString((*record)["last_event_datetime"].(string))
	ryosei.CreatedAt = (*record)["created_at"].(string)
	ryosei.UpdatedAt = toNullString((*record)["updated_at"].(string))
	ryosei.SharingStatus = floatToInt((*record)["sharing_status"].(float64))

	return nil
}

func InsertRyoseis(db *sqlx.DB, ryoseis []*Ryosei) error {
	return nil
}

func UpdateRyoseis(db *sqlx.DB, ryoseis []*Ryosei) error {
	return nil
}

func GetUnsyncedRyoseisAsSqlInsert(db *sqlx.DB) (*string, error) {
	sql := ""
	return &sql, nil
}

func GetUnsyncedRyoseisAsSqlUpdate(db *sqlx.DB) (*string, error) {
	sql := ""
	return &sql, nil
}

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
			`INSERT INTO parcels(
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
			kana,
			romaji,
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
