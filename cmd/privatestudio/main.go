package main

import (
	"github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/database"
	"github.com/skruger/privatestudio/gtk4"
	"os"
)

func main() {
	db, err := database.NewSqliteDb("gtk.sqlite")
	if err != nil {
		log.Panic("unable to open sql database", err)
	}
	err = database.MigrateUp(db)
	if err != nil {
		log.Error("unable to migrate sql database", err)
		os.Exit(1)
	}

	gtk4.RunUI(db)

}
