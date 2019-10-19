package main

import (
	"os"
	"path"
	"runtime"
	"syscall"
	"time"

	"net/http"
	"os/signal"

	"hackathon.com/pyz/api"
	"hackathon.com/pyz/dbm"
	"hackathon.com/pyz/env"

	"github.com/go-chi/chi"
)

var (
	environ = env.GetEnvironment()
	logger  = env.GetLogger()
)

//go:generate sh -c "echo \"package env\n\nvar buildVersion = \\\"`git describe --tags`\\\" \" > env/env_version_generated.go"
func main() {
	// use all cpus
	runtime.GOMAXPROCS(runtime.NumCPU())

	// configure mysql datastore
	db, err := dbm.NewDatabaseManager(environ, environ.AppName)
	if err != nil {
		logger.Fatal(err)
	}

	// check mysql database time is in sync with wallet
	var mysqlTime time.Time
	if err = db.Get(&mysqlTime, "SELECT NOW()"); err != nil {
		logger.Fatal(err)
	}

	// warn if server time is not in sync with mysql
	var prepareTime = time.Now().In(environ.Location)
	if prepareTime.Sub(mysqlTime) > time.Second {
		logger.Fatal("System time not sync with mysql server", mysqlTime, prepareTime)
	}

	// configure url routes & routes are important right now
	// Important: middleware order is important
	router := chi.NewRouter()
	router.Use(api.LoggingMiddleware)
	router.Use(api.PanicMiddleware)

	// mount endpoints
	router.Mount("/api/", api.NewRouter(db))

	// mainly to avoid middlewares for prom handler
	mainRouter := chi.NewRouter()
	mainRouter.Mount("/", router)

	// allow access to image directory with main router
	workDir, _ := os.Getwd()
	filesDir := path.Join(workDir, "assets")
	api.FileServer(mainRouter, "/assets", http.Dir(filesDir))

	// gracefully exit on interrupt & terminate signals
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	go graceFulExitOnSignal(exitSignal, db)

	logger.Printf("Starting %s at %s", environ.AppName, environ.Port)

	logger.Fatal(http.ListenAndServe(
		environ.Port,
		mainRouter,
	))
}

func graceFulExitOnSignal(exitSignal chan os.Signal, db *dbm.DatabaseManager) {
	<-exitSignal

	logger.Info("Waiting for all go routines to finish")

	db.DB.Close()
	os.Exit(0)
}
