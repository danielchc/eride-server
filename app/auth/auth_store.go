package auth

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

type AuthStore struct {
	db *sql.DB
}

// // UserStore is an interface to store users
// type UserStoreInterface interface {
// 	// Save saves a user to the store
// 	Save(user *User) error
// 	// Find finds a user by username
// 	Find(username string) (*User, error)
// }

func NewUserStore(db *sql.DB) *AuthStore {
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

	// Insert the new user
	_, err = store.db.Exec(`
		INSERT INTO users (username, password)
		VALUES (?, ?)
	`, user.Username, user.HashedPassword)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}

	return nil
}

// // Find finds a user by username
func (store *AuthStore) Find(username string) (*User, error) {
	row := store.db.QueryRow(`
		SELECT username, password
		FROM users
		WHERE username = ?
	`, username)

	var user User
	if err := row.Scan(&user.Username, &user.HashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err // Some other error
	}

	return &user, nil
}

// ErrAlreadyExists is returned when a user already exists in the store
var ErrAlreadyExists = errors.New("user already exists")
