package dbm

import (
	"fmt"

	"hackathon.com/pyz/env"

	_ "github.com/go-sql-driver/mysql" // mysql import
	"github.com/jmoiron/sqlx"
)

var (
	environ = env.GetEnvironment()
	logger  = env.GetLogger()
)

// DatabaseOperator interface define common database operations
type DatabaseOperator interface {
	ProfileOperations
	sqlx.Ext

	Beginx() (*sqlx.Tx, error)
}

// DatabaseManager implements DatabaseOperator interface
type DatabaseManager struct {
	*sqlx.DB
	*ProfileManager
}

// NewDatabaseManager return a new MySQLStore
func NewDatabaseManager(environ *env.Envs, prefix string) (*DatabaseManager, error) {
	var defaultDB *sqlx.DB
	var err error

	// connect to database and fail if any error
	if defaultDB, err = sqlx.Connect("mysql", environ.DatabaseURL); err != nil {
		return nil, err
	}

	// don't re-use idle connections after environ.DatabaseConnLifetime
	defaultDB.SetConnMaxLifetime(environ.DatabaseConnLifetime)

	var pm = NewProfileManager(defaultDB, prefix)

	return &DatabaseManager{
		DB:             defaultDB,
		ProfileManager: pm,
	}, nil
}

// NewTestDatabaseManager return a db manager for testing
func NewTestDatabaseManager() *DatabaseManager {
	var prefix = "test_" + environ.AppName

	if dbManager, err := NewDatabaseManager(environ, prefix); err != nil {
		panic(err)
	} else {
		return dbManager
	}
}

// CreateTables create all required tables in database for dbm
func (dbm DatabaseManager) CreateTables() {
	var createProfileSQL = fmt.Sprintf(createProfileSQL, dbm.ProfileManager.tableName, environ.DatabaseEngine)
	if _, err := dbm.Exec(createProfileSQL); err != nil {
		logger.Error(err)
	}

	var createActivePacksSQL = fmt.Sprintf(createActivePacksSQL, dbm.ProfileManager.activePackTableName, environ.DatabaseEngine)
	if _, err := dbm.Exec(createActivePacksSQL); err != nil {
		logger.Error(err)
	}
}

// DropTables drop all account and client tables from database
func (dbm DatabaseManager) DropTables() {
	dbm.MustExec(`DROP TABLE IF EXISTS ` + dbm.ProfileManager.tableName)
	dbm.MustExec(`DROP TABLE IF EXISTS ` + dbm.ProfileManager.activePackTableName)
}
