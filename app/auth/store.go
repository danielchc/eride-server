package auth

import (
	"chenel/eride/app/dto"
	"errors"
	"fmt"

	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type AuthStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// Save saves a user to the SQLite database
func (store *AuthStore) Save(user *dto.User) error {
	// Check if the user already exists
	existingUser, err := store.Find(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrAlreadyExists
	}

	// Insert the new user
	result := store.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to save user: %v", result.Error)
	}

	return nil
}

// Find finds a user by username
func (store *AuthStore) Find(username string) (*dto.User, error) {
	var user dto.User
	result := store.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // User not found
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// ErrAlreadyExists is returned when a user already exists in the store
var ErrAlreadyExists = errors.New("user already exists")
