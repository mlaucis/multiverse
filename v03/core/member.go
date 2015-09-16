package core

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

// Member interface
type Member interface {
	// Create adds a new account user to the database and returns the created account user or an error
	Create(accountUser *entity.Member, retrieve bool) (*entity.Member, []errors.Error)

	// Read returns the account matching the ID or an error
	Read(accountID, accountUserID int64) (accountUser *entity.Member, er []errors.Error)

	// Update update an account user in the database and returns the updated account user or an error
	Update(existingAccountUser, updatedAccountUser entity.Member, retrieve bool) (*entity.Member, []errors.Error)

	// Delete deletes the account user matching the IDs or an error
	Delete(*entity.Member) []errors.Error

	// List returns all the users from a certain account
	List(accountID int64) (accountUsers []*entity.Member, er []errors.Error)

	// CreateSession handles the creation of a user session and returns the session token
	CreateSession(user *entity.Member) (string, []errors.Error)

	// RefreshSession generates a new session token for the user session
	RefreshSession(sessionToken string, user *entity.Member) (string, []errors.Error)

	// DestroySession removes the user session
	DestroySession(sessionToken string, user *entity.Member) []errors.Error

	// GetSession retrieves the account user session token
	GetSession(user *entity.Member) (string, []errors.Error)

	// FindByEmail returns the account and account user for a certain e-mail address
	FindByEmail(email string) (*entity.Organization, *entity.Member, []errors.Error)

	// ExistsByEmail checks if the account exists for a certain email
	ExistsByEmail(email string) (bool, []errors.Error)

	// FindByUsername returns the account and account user for a certain username
	FindByUsername(username string) (*entity.Organization, *entity.Member, []errors.Error)

	// ExistsByUsername checks if the account exists for a certain username
	ExistsByUsername(username string) (bool, []errors.Error)

	// ExistsByID checks if an account user exists by ID or not
	ExistsByID(accountID, accountUserID int64) (bool, []errors.Error)

	// FindBySession will load an account by the session key, if it exists
	FindBySession(sessionKey string) (*entity.Member, []errors.Error)

	// FindByPublicID will load an account by the public ID it exposes to the world
	FindByPublicID(accountID int64, publicID string) (*entity.Member, []errors.Error)
}
