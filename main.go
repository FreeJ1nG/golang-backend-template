package main

import (
	"log"

	"github.com/FreeJ1nG/backend-template/app"
	"github.com/FreeJ1nG/backend-template/db"
	"github.com/FreeJ1nG/backend-template/util"
)

func main() {
	config, err := util.SetConfig()
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err.Error())
		return
	}

	mainDB := db.CreatePool(config.DBDsn)
	db.TestConnection(mainDB)
	defer func() {
		mainDB.Close()
	}()

	s := app.NewServer(config, mainDB)
	s.InjectDependencies()
	s.ListenAndServe()
}
