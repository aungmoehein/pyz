package dbm

import (
	"github.com/jmoiron/sqlx"
	"hackathon.com/pyz/env"
)

// ProfileOperations specify operations for profile
type ProfileOperations interface {
	GetProfile(int) (Profile, error)
	AddProfile(sqlx.Ext, Profile) error
	GetLeaderboard() ([]LeaderBoard, error)

	AddPack(int, string) error
	GetPacks(int) ([]env.Pack, error)
	HasActivePack(int) (bool, error)
}

// Profile represents a record in profile database
type Profile struct {
	ID        int    `db:"id" json:"-"`
	UID       int    `db:"uid" json:"uid"`
	Name      string `db:"name" json:"name"`
	Phone     string `db:"phone" json:"phone"`
	ImgURL    string `db:"imgurl" json:"img_url"`
	Active    bool   `db:"active" json:"active"`
	Score     int    `db:"score" json:"score"`
	CreatedAt Time   `db:"created_at" json:"created_at"`
	UpdatedAt Time   `db:"updated_at" json:"updated_at"`
}

// LeaderBoard represents a record in profile database
type LeaderBoard struct {
	UID   int    `db:"uid" json:"uid"`
	Name  string `db:"name" json:"name"`
	Score int    `db:"score" json:"score"`
	Rank  int    `db:"ranking" json:"ranking"`
}

// ActivePack contains activepack record from database
type ActivePack struct {
	ID        int    `db:"id" json:"-"`
	UID       int    `db:"uid" json:"uid"`
	PID       string `db:"pid" json:"pid"`
	Active    bool   `db:"active" json:"active"`
	CreatedAt Time   `db:"created_at" json:"created_at"`
	UpdatedAt Time   `db:"updated_at" json:"updated_at"`
}
