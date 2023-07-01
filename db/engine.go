package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseEngine struct {
	DB     *gorm.DB
	models map[string]interface{}
}

var engine *DatabaseEngine = &DatabaseEngine{}

func GetDatabaseEngine() (*DatabaseEngine, error) {
	if engine != nil && engine.DB != nil {
		return engine, nil
	}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 构建项目路径
	projectPath := filepath.Join(dir, "./") // 假设项目的根目录在当前目录的上一级目录

	// 构建SQLite数据库文件路径
	dbPath := filepath.Join(projectPath, "db/sqlite3.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error occurred during open sqlite database: %w", err)
	}

	engine.DB = db

	engine.models = map[string]interface{}{
		"host_info": &Host{},
		"share":     &Share{},
		"directory": &Directory{},
	}

	return engine, nil
}

func (engine *DatabaseEngine) Migrate() error {
	models := make([]interface{}, 0, len(engine.models))
	for _, model := range engine.models {
		models = append(models, model)
	}

	if err := engine.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("error occurred during auto migration: %w", err)
	}

	return nil
}
