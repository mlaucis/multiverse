package core

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/entity"
)

// Organization interface
type Organization interface {
	// Create adds a new organization to the database and returns the created account or an error
	Create(account *entity.Organization, retrieve bool) (*entity.Organization, []errors.Error)

	// Read returns the account matching the ID or an error
	Read(accountID int64) (*entity.Organization, []errors.Error)

	// Update updates the account matching the ID or an error
	Update(existingAccount, updatedAccount entity.Organization, retrieve bool) (*entity.Organization, []errors.Error)

	// Delete deletes the account matching the ID or an error
	Delete(*entity.Organization) []errors.Error

	// Exists validates if an account exists and returns the account or an error
	Exists(accountID int64) (bool, []errors.Error)

	// FindByKey will load an account by the account key, if it exists
	FindByKey(authKey string) (*entity.Organization, []errors.Error)

	// ReadByPublicID returns the account matching the public ID
	ReadByPublicID(id string) (*entity.Organization, []errors.Error)
}
