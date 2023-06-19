package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseEngine struct {
	DB     *gorm.DB
	Models map[string]interface{}
}

var engine *DatabaseEngine = &DatabaseEngine{}

func GetDatabaseEngine() (*DatabaseEngine, error) {
	if engine != nil && engine.DB != nil {
		return engine, nil
	}

	db, err := gorm.Open(sqlite.Open("./db/sqlite3.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error occurred during open sqlite database: %w", err)
	}

	engine.DB = db

	engine.Models = map[string]interface{}{
		"host_info": &HostInfo{},
		"share":     &Share{},
	}

	return engine, nil
}

func (engine *DatabaseEngine) Migrate() error {
	models := make([]interface{}, 0, len(engine.Models))
	for _, model := range engine.Models {
		models = append(models, model)
	}

	if err := engine.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("error occurred during auto migration: %w", err)
	}

	return nil
}

// func (engine *DatabaseEngine) GetModel(modelName string) (interface{}, error) {
// 	model, ok := engine.Models[modelName]
// 	if !ok {
// 		return nil, fmt.Errorf("table '%s' does not exist", modelName)
// 	}
// 	return model, nil
// }
