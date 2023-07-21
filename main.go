package main

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/webservice"
)

func main() {
	common.ConfigLogger()

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	webservice.Start()
}
