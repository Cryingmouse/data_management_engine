package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cryingmouse/data_management_engine/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var engine *DatabaseEngine

// DatabaseEngine struct holds the database connection and models.
type DatabaseEngine struct {
	DB     *gorm.DB
	Models map[string]interface{}
}

// GetDatabaseEngine returns the instance of DatabaseEngine.
func GetDatabaseEngine() (*DatabaseEngine, error) {
	if engine != nil && engine.DB != nil {
		return engine, nil
	}

	// Get the current working directory. Assuming the root directory of the project is the current directory
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the SQLite database file path
	dbPath := filepath.Join(dir, "db/sqlite3.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: common.DBLogger})
	if err != nil {
		return nil, fmt.Errorf("error occurred while opening SQLite database: %w", err)
	}

	engine = &DatabaseEngine{
		DB: db.Debug(),
		Models: map[string]interface{}{
			"host_info":  &Host{},
			"share":      &CIFSShare{},
			"directory":  &Directory{},
			"local_user": &LocalUser{},
		},
	}

	return engine, nil
}

// Migrate performs auto migration for all registered models.
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
