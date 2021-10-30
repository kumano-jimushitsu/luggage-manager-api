package models

import "database/sql"

type Ryosei struct {
	Id                int
	Room              string
	Name              string
	Kana              string
	Romaji            string
	BlockID           int
	SlackID           string
	Status            int
	CurrentCount      int
	TotalCount        int
	TotalWait         string
	LastEventID       int
	LastEventDatetime string
	CreatedAt         string
	UpdatedAt         string
	SharingStatus     string
}

func AllRyoseis(db *sql.DB) ([]*Ryosei, error) {
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
			&ryosei.Room,
			&ryosei.Name,
			&ryosei.Kana,
			&ryosei.Romaji,
			&ryosei.BlockID,
			&ryosei.SlackID,
			&ryosei.Status,
			&ryosei.CurrentCount,
			&ryosei.TotalCount,
			&ryosei.TotalWait,
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
