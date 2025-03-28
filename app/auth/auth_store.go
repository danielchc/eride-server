package auth

import (
	"chenel/eride/app/db"
	"errors"

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
func (store *AuthStore) Save(user *User) error {
	// Check if the user already exists
	existingUser, err := store.Find(user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrAlreadyExists
	}

	//???????

	store.db.Create(&db.User{FirstName: "Alice", LastName: "Smith", Username: user.Username, Password: user.HashedPassword, TOTPSecret: "TOTP_SECRET_1"})

	// Insert the new user
	// _, err = store.db.
	// if err != nil {
	// 	return fmt.Errorf("failed to save user: %v", err)
	// }

	return nil
}

// // Find finds a user by username
func (store *AuthStore) Find(username string) (*User, error) {
	var user User
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
