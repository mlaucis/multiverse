/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// AccountUser interface
	AccountUser interface {
		// Create adds a new account user to the database and returns the created account user or an error
		Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error)

		// Read returns the account matching the ID or an error
		Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er errors.Error)

		// Update update an account user in the database and returns the updated account user or an error
		Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error)

		// Delete deletes the account user matching the IDs or an error
		Delete(*entity.AccountUser) errors.Error

		// List returns all the users from a certain account
		List(accountID int64) (accountUsers []*entity.AccountUser, er errors.Error)

		// CreateSession handles the creation of a user session and returns the session token
		CreateSession(user *entity.AccountUser) (string, errors.Error)

		// RefreshSession generates a new session token for the user session
		RefreshSession(sessionToken string, user *entity.AccountUser) (string, errors.Error)

		// DestroySession removes the user session
		DestroySession(sessionToken string, user *entity.AccountUser) errors.Error

		// GetSession retrieves the account user session token
		GetSession(user *entity.AccountUser) (string, errors.Error)

		// FindByEmail returns the account and account user for a certain e-mail address
		FindByEmail(email string) (*entity.Account, *entity.AccountUser, errors.Error)

		// ExistsByEmail checks if the account exists for a certain email
		ExistsByEmail(email string) (bool, errors.Error)

		// FindByUsername returns the account and account user for a certain username
		FindByUsername(username string) (*entity.Account, *entity.AccountUser, errors.Error)

		// ExistsByUsername checks if the account exists for a certain username
		ExistsByUsername(username string) (bool, errors.Error)

		// ExistsByID checks if an account user exists by ID or not
		ExistsByID(accountID, accountUserID int64) (bool, errors.Error)

		// FindBySession will load an account by the session key, if it exists
		FindBySession(sessionKey string) (*entity.AccountUser, errors.Error)

		// FindByPublicID will load an account by the public ID it exposes to the world
		FindByPublicID(accountID int64, publicID string) (*entity.AccountUser, errors.Error)
	}
)
