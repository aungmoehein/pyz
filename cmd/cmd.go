package main

import (
	"github.com/trhura/simplecli"

	"hackathon.com/pyz/env"
)

var environ = env.GetEnvironment()
var logger = env.GetLogger()

// CLIOperations contains all CLI related operations
type CLIOperations struct {
	DB *DatabaseOperations
}

func main() {
	simplecli.Handle(&CLIOperations{
		DB: &DatabaseOperations{},
	})
}
