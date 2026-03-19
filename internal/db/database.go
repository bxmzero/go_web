package db

import (
	"go_web/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("demo.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := database.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return database, nil
}
