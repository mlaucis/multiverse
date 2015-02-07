/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package mysql

import (
	"fmt"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core/entity"
)

// GetApplicationUserByToken returns the user corresponding to the applicationId / userToken combination or an error
func GetApplicationUserByToken(applicationId uint64, userID string) (user *entity.User, err error) {
	user = &entity.User{}

	// Execute query to get user
	err = GetSlave().
		QueryRowx("SELECT * FROM `users` WHERE `application_id`=? AND `token`=?", applicationId, userID).
		StructScan(user)

	return
}

// GetApplicationUsers returns all the users corresponding to the applicationId / userToken combination or an error
func GetApplicationUsers(applicationId uint64) (users []*entity.User, err error) {
	users = []*entity.User{}

	// Execute query to get users
	err = GetSlave().
		Select(&users, "SELECT * FROM `users` WHERE `application_id`=?", applicationId)

	return
}

// AddApplicationUser creates a user for an account and returns the created entry or an error
func AddApplicationUser(applicationId uint64, user *entity.User) (*entity.User, error) {
	// Check if token empty
	if user.AuthToken == "" {
		return nil, fmt.Errorf("empty user token is not allowed")
	}
	// Check if token empty
	if user.Username == "" {
		return nil, fmt.Errorf("empty user username is not allowed")
	}
	// Check if token empty
	if user.Password == "" {
		return nil, fmt.Errorf("empty user password is not allowed")
	}

	// Write to db
	query := "INSERT INTO `users` (`application_id`, `token`, `username`, `name`, `password`, `email`, `url`, `thumbnail_url`, `custom`)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := GetMaster().
		Exec(query, applicationId, user.AuthToken, user.Username, user.FirstName, user.Password, user.Email, user.URL, user.Image[0].URL, user.Metadata)
	if err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return nil, fmt.Errorf("error while saving to database")
	}

	// Return application user
	return GetApplicationUserByToken(applicationId, user.AuthToken)
}

// GetApplicationUserWithConnections returns the user and it's connections to other users or an error
func GetApplicationUserWithConnections(applicationId uint64, userID string) (user *entity.User, err error) {
	// Check if token empty
	if userID == "" {
		return nil, fmt.Errorf("empty user token is not allowed")
	}
	// Retrieve user
	if user, err = GetApplicationUserByToken(applicationId, userID); err != nil {
		return
	}
	user.Connections = []*entity.User{}

	// Retrieve user connections
	err = GetSlave().
		Select(
		&user.Connections,
		"SELECT `users`.* "+
			"FROM `users` "+
			"JOIN `user_connections` as `guc` "+
			"ON `users`.`application_id` = `guc`.`application_id` AND "+
			"`users`.`token` = `guc`.`user_id2` "+
			"WHERE `guc`.`application_id`=? AND `guc`.`user_id1`=?",
		applicationId,
		userID,
	)

	return
}

// AddApplicationUserConnection will add a new connection between users or returns an error
func AddApplicationUserConnection(applicationId uint64, connection *entity.Connection) (err error) {
	// Check if token1 empty
	if connection.UserFromID == 0 {
		return fmt.Errorf("empty user1 token is not allowed")
	}
	// Check if token2 empty
	if connection.UserToID == 0 {
		return fmt.Errorf("empty user2 token is not allowed")
	}
	// Write to db
	query := "INSERT INTO `user_connections` (`application_id`, `user_id1`, `user_id2`) VALUES (?, ?, ?)"
	_, err = GetMaster().
		Exec(query, applicationId, connection.UserFromID, connection.UserToID)
	if err != nil {
		if config.Conf().Env() == "dev" {
			return err
		}
		return fmt.Errorf("error while saving to database")
	}

	return
}
