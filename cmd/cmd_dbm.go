package main

import (
	"log"

	"hackathon.com/pyz/dbm"
)

// DatabaseOperations for database operations
type DatabaseOperations struct{}

// Create required tables for the app
func (h DatabaseOperations) Create() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()

	var err error
	var manager *dbm.DatabaseManager
	if manager, err = dbm.NewDatabaseManager(environ, environ.AppName); err != nil {
		log.Fatal("Cannot open database " + environ.DatabaseURL)
	}

	manager.CreateTables()
}
