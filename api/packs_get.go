package api

import (
	"hackathon.com/pyz/env"
	"net/http"
)

// PacksHandler get news request
type PacksHandler Handler

func (h PacksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid int
	var packs []env.Pack

	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	if _, err = h.DB.GetProfile(uid); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	for _, service := range environ.Packs.Packs {
		packs = append(packs, service)
	}

	WriteJSON(w, PacksResponse{
		Success: true,
		Packs:   packs,
	})
}
