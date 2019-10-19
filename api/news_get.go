package api

import (
	"hackathon.com/pyz/env"
	"net/http"
)

// NewsHandler get news request
type NewsHandler Handler

func (h NewsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var uid int
	var news []env.New

	if uid, err = GetUID(r); err != nil {
		panic(ErrInvalidUID.Wraps(err))
	}

	if _, err = h.DB.GetProfile(uid); err != nil {
		panic(ErrDatabase.Wraps(err))
	}

	for _, article := range environ.Article.News {
		news = append(news, article)
	}

	WriteJSON(w, NewsResponse{
		Success: true,
		News:    news,
	})
}
