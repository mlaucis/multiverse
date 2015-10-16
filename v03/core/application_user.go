package core

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

// ApplicationUser interface
type ApplicationUser interface {
	// Create adds a user to the database and returns the created user or an error
	Create(accountID, applicationID int64, user *entity.ApplicationUser) []errors.Error

	// Read returns the user matching the ID or an error
	Read(accountID, applicationID int64, userID uint64, withStatistics bool) (user *entity.ApplicationUser, err []errors.Error)

	// ReadMultiple returns all users matching the desired IDs
	ReadMultiple(accountID, applicationID int64, userIDs []uint64) (users []*entity.ApplicationUser, err []errors.Error)

	// Update updates a user in the database and returns the updates user or an error
	Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err []errors.Error)

	// Delete deletes the user matching the IDs or an error
	Delete(accountID, applicationID int64, userID uint64) (err []errors.Error)

	// List returns all users from a certain account
	List(accountID, applicationID int64) (users []*entity.ApplicationUser, err []errors.Error)

	// CreateSession handles the creation of a user session and returns the session token
	CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error)

	// RefreshSession generates a new session token for the user session
	RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, []errors.Error)

	// GetSession returns the application user session
	GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error)

	// DestroySession removes the user session
	DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) []errors.Error

	// FindByEmail returns an application user by its email
	FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, []errors.Error)

	// ExistsByEmail checks if an application user exists by searching it via the email
	ExistsByEmail(accountID, applicationID int64, email string) (bool, []errors.Error)

	// FindByUsername returns an application user by its username
	FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, []errors.Error)

	// ExistsByUsername checks if an application user exists by searching it via the username
	ExistsByUsername(accountID, applicationID int64, username string) (bool, []errors.Error)

	// ExistsByID validates if a user exists and returns it or an error
	ExistsByID(accountID, applicationID int64, userID uint64) (bool, []errors.Error)

	// FindBySession will load an application user by the session key, if it exists
	FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, []errors.Error)

	// Search finds the users matching the search term
	Search(accountID, applicationID int64, searchTerm string) (user []*entity.ApplicationUser, err []errors.Error)

	// Read friend and follower statistics for a user
	FriendStatistics(accountID, applicationID int64, appUser *entity.ApplicationUser) []errors.Error
}
