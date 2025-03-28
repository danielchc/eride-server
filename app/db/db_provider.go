package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID         uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	TOTPSecret string `gorm:"not null"`
	CreatedAt  time.Time
}

// Group model
type Group struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// Vault model
type Vault struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Desc      string
	CreatedAt time.Time
	Folders   []Folder
}

// Folder model
type Folder struct {
	ID       uint `gorm:"primaryKey"`
	VaultID  uint `gorm:"not null"`
	ParentID *uint
	Name     string `gorm:"not null"`
	Entries  []Entry
}

// Entry model
type Entry struct {
	ID       uint   `gorm:"primaryKey"`
	FolderID uint   `gorm:"not null"`
	Name     string `gorm:"not null"`
}

// EntryRevision model
type EntryRevision struct {
	ID        uint `gorm:"primaryKey"`
	EntryID   uint `gorm:"not null"`
	Version   int  `gorm:"not null"`
	CreatedAt time.Time
}

// ACL model
type ACL struct {
	ID      uint `gorm:"primaryKey"`
	VaultID uint `gorm:"not null"`
	UserID  *uint
	GroupID *uint
	Role    string `gorm:"not null"`
}

// Tag model
type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// EntryTag model
type EntryTag struct {
	EntryID uint `gorm:"not null"`
	TagID   uint `gorm:"not null"`
}

func NewDBProvider(dbPath string) (*gorm.DB, error) {
	// Open database
	db, err := gorm.Open(sqlite.Open("password_manager.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	// Auto migrate schemas
	err = db.AutoMigrate(&User{}, &Group{}, &Vault{}, &Folder{}, &Entry{}, &EntryRevision{}, &ACL{}, &Tag{}, &EntryTag{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database file: %v", err)
	}

	return db, nil

}
