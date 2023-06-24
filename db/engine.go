package db

import (
	"fmt"

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

	db, err := gorm.Open(sqlite.Open("./sqlite3.db"), &gorm.Config{})
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

func (engine *DatabaseEngine) Get(dest interface{}) *gorm.DB {

	return engine.DB.Where(dest).First(dest)
}
