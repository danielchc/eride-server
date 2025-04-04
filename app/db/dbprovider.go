package db

import (
	"fmt"
	"log"

	"chenel/eride/app/dto"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDBProvider(dbPath string) (*gorm.DB, error) {
	return NewDBProvider(sqlite.Open(dbPath))
}

func NewDBProvider(dialector gorm.Dialector) (*gorm.DB, error) {
	// Open database
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto migrate schemas
	err = db.AutoMigrate(&dto.User{}, &dto.Group{}, &dto.Vault{}, &dto.Folder{}, &dto.Entry{}, &dto.EntryRevision{}, &dto.ACL{}, &dto.Tag{}, &dto.EntryTag{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database file: %v", err)
	}

	return db, nil

}
