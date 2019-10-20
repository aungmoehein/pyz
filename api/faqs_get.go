package api

import (
	"hackathon.com/pyz/env"
	"net/http"
)

// FaqsHandler get news request
type FaqsHandler Handler

func (h FaqsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid int
	var faqs []env.FAQ

	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	if _, err = h.DB.GetProfile(uid); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	for _, faq := range environ.FAQs.FAQs {
		faqs = append(faqs, faq)
	}

	WriteJSON(w, FaqsResponse{
		Success: true,
		FAQ:     faqs,
	})
}
