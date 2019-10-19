package api

import (
	"net/http"
	"strings"

	"hackathon.com/pyz/dbm"

	"github.com/jmoiron/sqlx"
)

// ProfileRegisterJSON contains json parameters for ProfileRegisterHandler
type ProfileRegisterJSON struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Active *bool  `json:"active"`
}

// ProfileRegisterHandler add profile request
type ProfileRegisterHandler Handler

func (h ProfileRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid int

	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	var jsonParams ProfileRegisterJSON
	if err = ReadJSON(r, &jsonParams); err != nil {
		panic(ErrRequestParams.Wraps(err))
	}

	var imgURL = ""
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if imgURL, err = SaveFile(r, "/profile/", uid, "_profile", ".png"); err != nil {
			panic(ErrUnexpected.Wraps(err))
		}
	}

	var tx *sqlx.Tx
	if tx, err = h.DB.Beginx(); err != nil {
		panic(ErrDatabaseTx.Wraps(err))
	}

	defer RollbackOnPanic(tx)

	var profile dbm.Profile
	if profile, err = h.DB.GetProfile(uid); err != nil {
		if err != dbm.ErrNoRowAffected {
			panic(ErrDatabase.Wraps(err))
		}

		// New profile, initialize default values
		profile.UID = uid
		profile.Active = true
	}

	if jsonParams.Active != nil {
		profile.Active = *jsonParams.Active
	}

	if jsonParams.Name != "" {
		profile.Name = jsonParams.Name
	}

	if jsonParams.Phone != "" {
		profile.Phone = jsonParams.Phone
	}

	if imgURL != "" {
		profile.ImgURL = imgURL
	}

	profile = dbm.Profile{
		UID:    profile.UID,
		Active: profile.Active,
		Name:   profile.Name,
		Phone:  profile.Phone,
		ImgURL: profile.ImgURL,
	}

	if err = h.DB.AddProfile(tx, profile); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	if err = tx.Commit(); err != nil {
		panic(ErrDatabaseTx.Wraps(err))
	}

	WriteJSON(w, DefaultResponse{
		Success: true,
		Message: "Profile added successfully",
	})
}
