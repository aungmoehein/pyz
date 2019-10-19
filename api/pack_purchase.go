package api

import (
	"net/http"
)

// PackPurchaseJSON contains json parameters for PackPurchaseJSON
type PackPurchaseJSON struct {
	PID string `json:"pid"`
}

// PackPurchaseHandler add pack purchase request
type PackPurchaseHandler Handler

func (h PackPurchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid int

	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	var jsonParams PackPurchaseJSON
	if err = ReadJSON(r, &jsonParams); err != nil {
		panic(ErrRequestParams.Wraps(err))
	}

	if err = h.DB.AddPack(uid, jsonParams.PID); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	WriteJSON(w, DefaultResponse{
		Success: true,
		Message: "Pack added successfully",
	})
}
