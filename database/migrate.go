package database

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/labstack/gommon/log"
)

//go:embed migrations/*.sql
var fs embed.FS

func MigrateUp(db *sql.DB) error {
	driver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	dbInstance, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance(
		"iofs",
		driver,
		"sqlite",
		dbInstance)
	if err != nil {
		return err
	}

	version, dirty, err := migration.Version()
	if version > 0 && dirty {
		log.Infof("Migration marked dirty with version %i and error %s", version, err)
	}
	if err != nil {
		log.Infof("unable to check migration version: %s", err)
	}
	if version < 1 {
		return migration.Up()
	}
	return nil
}
