package dbm

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"hackathon.com/pyz/env"
)

// ProfileManager handles profile operations
type ProfileManager struct {
	*sqlx.DB
	tableName           string
	activePackTableName string
}

// NewProfileManager returns new profile manager
func NewProfileManager(db *sqlx.DB, prefix string) *ProfileManager {
	return &ProfileManager{
		DB:                  db,
		tableName:           prefix + "_profile",
		activePackTableName: prefix + "_activepack",
	}
}

// AddProfile add a new profile record
func (um *ProfileManager) AddProfile(exec sqlx.Ext, profile Profile) error {
	var stmt = fmt.Sprintf(`
		INSERT INTO %s (uid, name, phone, imgurl, active, score) 
			VALUES (?, ?, ?, ?, ?, ?) 
			ON DUPLICATE KEY UPDATE name = VALUES(name),
			phone = VALUES(phone),
			imgurl = VALUES(imgurl),
			active = VALUES(active),
			score = VALUES(score)
	`, um.tableName)

	var args = make([]interface{}, 6)
	var err error

	args[0] = profile.UID
	args[1] = profile.Name
	args[2] = profile.Phone
	args[3] = profile.ImgURL
	args[4] = profile.Active
	args[5] = profile.Score

	if _, err = exec.Exec(stmt, args...); err != nil {
		return err
	}

	return nil
}

// GetProfile returns a new profile from database
func (um *ProfileManager) GetProfile(uid int) (Profile, error) {
	var profile Profile
	var err error

	var stmt = fmt.Sprintf(`
		SELECT * FROM %s WHERE uid = ?`, um.tableName)

	if err = um.Get(&profile, stmt, uid); err != nil {
		return Profile{}, err
	}

	return profile, nil
}

// GetLeaderboard return leaderboard list
func (um *ProfileManager) GetLeaderboard() ([]LeaderBoard, error) {
	var leaderboards []LeaderBoard
	var err error

	var stmt = fmt.Sprintf(`
		SELECT uid, name, score, FIND_IN_SET( 
		score, (SELECT GROUP_CONCAT( score
		ORDER BY score DESC ) 
		FROM %s )
		) AS ranking
		FROM %s ORDER BY ranking ASC`, um.tableName, um.tableName)

	if err = um.Select(&leaderboards, stmt); err != nil {
		return nil, err
	}

	return leaderboards, err
}

// AddPack add a new profile record
func (um *ProfileManager) AddPack(uid int, pid string) error {
	var stmt = fmt.Sprintf(`
		INSERT INTO %s (uid, pid, active) 
			VALUES (?, ?, ?)
	`, um.activePackTableName)

	var args = make([]interface{}, 3)
	var err error

	args[0] = uid
	args[1] = pid
	args[2] = true

	if _, err = um.Exec(stmt, args...); err != nil {
		return err
	}

	return nil
}

// GetPacks return active packs
func (um *ProfileManager) GetPacks(uid int) ([]env.Pack, error) {
	var activePacks []ActivePack
	var err error

	var stmt = fmt.Sprintf(`
		SELECT * FROM %s WHERE uid = ?`, um.activePackTableName)

	if err = um.Select(&activePacks, stmt, uid); err != nil {
		return nil, err
	}

	var packs []env.Pack
	for _, activePack := range activePacks {
		packs = append(packs, environ.Packs.Packs[activePack.PID])
	}

	return packs, nil
}

// HasActivePack check active status record
func (um *ProfileManager) HasActivePack(uid int) (bool, error) {
	var activePack ActivePack
	var err error

	var stmt = fmt.Sprintf(`
		SELECT * FROM %s WHERE uid = ? AND active = true`, um.activePackTableName)

	if err = um.Get(&activePack, stmt, uid); err != nil {
		return false, err
	}

	if activePack != (ActivePack{}) {
		return true, err
	}

	return false, nil
}
