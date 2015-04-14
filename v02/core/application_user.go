/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// ApplicationUser interface
	ApplicationUser interface {
		// Create adds a user to the database and returns the created user or an error
		Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError)

		// Read returns the user matching the ID or an error
		Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err tgerrors.TGError)

		// Update updates a user in the database and returns the updates user or an error
		Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError)

		// Delete deletes the user matching the IDs or an error
		Delete(accountID, applicationID, userID int64) (err tgerrors.TGError)

		// List returns all users from a certain account
		List(accountID, applicationID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)

		// CreateSession handles the creation of a user session and returns the session token
		CreateSession(user *entity.ApplicationUser) (string, tgerrors.TGError)

		// RefreshSession generates a new session token for the user session
		RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, tgerrors.TGError)

		// GetSession returns the application user session
		GetSession(user *entity.ApplicationUser) (string, tgerrors.TGError)

		// DestroySession removes the user session
		DestroySession(sessionToken string, user *entity.ApplicationUser) tgerrors.TGError

		// FindByEmail returns an application user by its email
		FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, tgerrors.TGError)

		// ExistsByEmail checks if an application user exists by searching it via the email
		ExistsByEmail(accountID, applicationID int64, email string) (bool, tgerrors.TGError)

		// FindByUsername returns an application user by its username
		FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, tgerrors.TGError)

		// ExistsByUsername checks if an application user exists by searching it via the username
		ExistsByUsername(accountID, applicationID int64, email string) (bool, tgerrors.TGError)

		// ExistsByID validates if a user exists and returns it or an error
		ExistsByID(accountID, applicationID, userID int64) bool
	}
)
