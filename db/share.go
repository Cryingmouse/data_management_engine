package db

import (
	"gorm.io/gorm"
)

type Share struct {
	gorm.Model
	Name        string `gorm:"unique"`
}
