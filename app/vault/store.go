package vault

import (
	"chenel/eride/app/dto"
	"fmt"

	"gorm.io/gorm"
)

type VaultStore struct {
	db *gorm.DB
}

func NewVaultStore(db *gorm.DB) *VaultStore {
	return &VaultStore{
		db: db,
	}
}

func (store *VaultStore) Save(vault *dto.Vault) error {
	result := store.db.Create(vault)
	if result.Error != nil {
		return fmt.Errorf("failed to save vault: %v", result.Error)
	}

	if result.Error != nil {
		return fmt.Errorf("failed to save vault: %v", result.Error)
	}

	return nil
}
