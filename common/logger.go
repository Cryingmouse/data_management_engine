package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func ConfigLogger() {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})

	// 打开日志文件
	file, err := os.OpenFile(config.Logger.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Println("Failed to open log file:", err)
	}

	// 设置日志级别
	switch config.Logger.LogLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}
}
