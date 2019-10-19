package api

import (
	"net/http"

	"hackathon.com/pyz/dbm"
)

// ProfileHandler get profile request
type ProfileHandler Handler

func (h ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	var uid int
	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	var profile dbm.Profile
	if profile, err = h.DB.GetProfile(uid); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	WriteJSON(w, ProfileResponse{
		Success: true,
		Profile: profile,
	})
}
