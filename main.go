package main

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/webservice"
)

func init() {
	// Access the Port configuration and assign it to the submodule's global variable
	if err := common.InitializeConfig("config.ini"); err != nil {
		panic(err)
	}
}

func main() {
	common.InitializeLogger()

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	webservice.Start()
}
