package api

import (
	// "net/http"
	// "time"

	"hackathon.com/pyz/dbm"
	"hackathon.com/pyz/env"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	environ  = env.GetEnvironment()
	logger   = env.GetLogger()
	decoder  = schema.NewDecoder()
	validate = validator.New()
)

// NewRouter returns with all routes
func NewRouter(dbm *dbm.DatabaseManager) chi.Router {
	var router = chi.NewRouter()
	router.Method("GET", "/user/{uid}/profile", &ProfileHandler{DB: dbm})
	router.Method("POST", "/user/{uid}/register", &ProfileRegisterHandler{DB: dbm})
	router.Method("GET", "/user/{uid}/news", &NewsHandler{DB: dbm})
	router.Method("GET", "/user/{uid}/packs", &PacksHandler{DB: dbm})
	router.Method("POST", "/user/{uid}/pack/purchase", &PackPurchaseHandler{DB: dbm})

	return router
}
