package util

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func ConnectDb(ConnStr string) error {
	var err error
	db, err = gorm.Open(postgres.Open(ConnStr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	return nil
}

func GetDb() *gorm.DB {
	return db
}
