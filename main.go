package main

import (
	"github.com/cryingmouse/data_management_engine/common"
	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/scheduler"
	"github.com/cryingmouse/data_management_engine/webservice"
)

func init() {
	// Initialize the configuration first and then Loggers.
	if err := common.InitializeConfig("config.ini"); err != nil {
		panic(err)
	}

	common.SetupLoggers()
	common.Logger.Debug("Initialize logger successfully.")
}

func main() {
	engine, err := db.GetDatabaseEngine()
	if err != nil {
		common.Logger.Error("Failed to initialize database. Error: %w", err)
		panic(err)
	}
	common.Logger.Debug("Initialize database successfully.")

	if err := engine.Migrate(); err != nil {
		common.Logger.Error("Failed to migration database. Error: %w", err)
	}

	scheduler.StartScheduler()

	webservice.Start()
	common.Logger.Debug("Start web service successfully.")
}
