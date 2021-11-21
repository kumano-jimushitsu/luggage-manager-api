package models

import "github.com/jmoiron/sqlx"

type Ryosei struct {
	Id                string `json:`
	RoomID            string `json:`
	Name              string `json:`
	Kana              string `json:`
	Romaji            string `json:`
	BlockID           int    `json:`
	SlackID           string `json:`
	Status            bool   `json:`
	CurrentCount      int    `json:`
	TotalCount        int    `json:`
	TotalWaitTime     string `json:`
	LastEventID       int    `json:`
	LastEventDatetime string `json:`
	CreatedAt         string `json:`
	UpdatedAt         string `json:`
	SharingStatus     string `json:`
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
