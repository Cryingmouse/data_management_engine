package main

import (
	"os"

	"github.com/cryingmouse/data_management_engine/db"
	"github.com/cryingmouse/data_management_engine/webservice"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// 打开日志文件
	file, err := os.OpenFile("dme.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Println("Failed to open log file:", err)
	}

	// 设置日志级别
	log.SetLevel(log.DebugLevel)

	engine, err := db.GetDatabaseEngine()
	if err != nil {
		panic(err)
	}

	engine.Migrate()

	webservice.Start()
}
