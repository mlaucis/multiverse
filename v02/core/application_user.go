/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// ApplicationUser interface
	ApplicationUser interface {
		// Create adds a user to the database and returns the created user or an error
		Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error)

		// Read returns the user matching the ID or an error
		Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err errors.Error)

		// Update updates a user in the database and returns the updates user or an error
		Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error)

		// Delete deletes the user matching the IDs or an error
		Delete(*entity.ApplicationUser) (err errors.Error)

		// List returns all users from a certain account
		List(accountID, applicationID int64) (users []*entity.ApplicationUser, err errors.Error)

		// CreateSession handles the creation of a user session and returns the session token
		CreateSession(user *entity.ApplicationUser) (string, errors.Error)

		// RefreshSession generates a new session token for the user session
		RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, errors.Error)

		// GetSession returns the application user session
		GetSession(user *entity.ApplicationUser) (string, errors.Error)

		// DestroySession removes the user session
		DestroySession(sessionToken string, user *entity.ApplicationUser) errors.Error

		// FindByEmail returns an application user by its email
		FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, errors.Error)

		// ExistsByEmail checks if an application user exists by searching it via the email
		ExistsByEmail(accountID, applicationID int64, email string) (bool, errors.Error)

		// FindByUsername returns an application user by its username
		FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, errors.Error)

		// ExistsByUsername checks if an application user exists by searching it via the username
		ExistsByUsername(accountID, applicationID int64, username string) (bool, errors.Error)

		// ExistsByID validates if a user exists and returns it or an error
		ExistsByID(accountID, applicationID, userID int64) (bool, errors.Error)

		// FindBySession will load an application user by the session key, if it exists
		FindBySession(sessionKey string) (*entity.ApplicationUser, errors.Error)
	}
)
