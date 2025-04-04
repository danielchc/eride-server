package dto

import "time"

// User model
type User struct {
	ID         uint64 `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	TOTPSecret string `gorm:"not null"`
	CreatedAt  time.Time
}

// Group model
type Group struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// Vault model
type Vault struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Desc      string
	CreatedAt time.Time
	Folders   []Folder
	ACL       []ACL `gorm:"foreignKey:VaultID"`
}

// ACL model
type ACL struct {
	ID      uint64 `gorm:"primaryKey"`
	UserID  *uint64
	GroupID *uint64
	Role    string `gorm:"not null"`
	VaultID uint64 `gorm:"not null"`
}

// Folder model
type Folder struct {
	ID       uint64 `gorm:"primaryKey"`
	VaultID  uint64 `gorm:"not null"`
	ParentID *uint64
	Name     string `gorm:"not null"`
	Entries  []Entry
}

// Entry model
type Entry struct {
	ID       uint64 `gorm:"primaryKey"`
	FolderID uint64 `gorm:"not null"`
	Name     string `gorm:"not null"`
}

// EntryRevision model
type EntryRevision struct {
	ID        uint64 `gorm:"primaryKey"`
	EntryID   uint64 `gorm:"not null"`
	Version   int    `gorm:"not null"`
	CreatedAt time.Time
}

// Tag model
type Tag struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// EntryTag model
type EntryTag struct {
	EntryID uint64 `gorm:"not null"`
	TagID   uint64 `gorm:"not null"`
}
