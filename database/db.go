package database

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
)

func NewSqliteDb(dbFile string) (*sql.DB, error) {
	dbConnectionString := fmt.Sprintf("file:%s?cache=shared", dbFile)
	log.Infof("Sqlite connection string: %s", dbConnectionString)
	return sql.Open("sqlite", dbConnectionString)
}
