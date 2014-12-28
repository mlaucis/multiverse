/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/gluee/backend/config"
	"github.com/gluee/backend/entity"
)

// GetApplicationUserByToken returns the user corresponding to the appID / userToken combination or an error
func GetApplicationUserByToken(appID uint64, userToken string) (user *entity.User, err error) {
	user = &entity.User{}

	err = GetSlave().
		QueryRowx("SELECT * FROM `gluee`.`users` WHERE `app_id`=? AND `token`=?", appID, userToken).
		StructScan(user)

	return
}

// GetApplicationUsers returns all the users corresponding to the appID / userToken combination or an error
func GetApplicationUsers(appID uint64) (users []*entity.User, err error) {
	users = []*entity.User{}

	err = GetSlave().
		Select(&users, "SELECT * FROM `gluee`.`users` WHERE `app_id`=?", appID)

	return
}

// AddApplicationUser creates a user for an account and returns the created entry or an error
func AddApplicationUser(appID uint64, user *entity.User) (*entity.User, error) {
	query := "INSERT INTO `gluee`.`users` (`app_id`, `token`, `username`, `name`, `password`, `email`, `url`, `thumbnail_url`, `custom`)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := GetMaster().
		Exec(query, appID, user.Token, user.Username, user.Name, user.Password, user.Email, user.URL, user.ThumbnailURL, user.Custom)
	if err != nil {
		if config.GetConfig().Env == "dev" {
			return nil, err
		}
		return nil, fmt.Errorf("error while saving to database")
	}

	return GetApplicationUserByToken(appID, user.Token)
}

// GetApplicationUserWithConnections returns the user and it's connections to other users or an error
func GetApplicationUserWithConnections(appID uint64, userToken string) (user *entity.User, err error) {
	if user, err = GetApplicationUserByToken(appID, userToken); err != nil {
		return
	}

	user.Connections = []*entity.User{}
	err = GetSlave().
		Select(
		&user.Connections,
		"SELECT `gluee`.`users`.* "+
			"FROM `gluee`.`users` "+
			"JOIN `gluee`.`user_connections` as `guc` "+
			"ON `gluee`.`users`.`app_id` = `guc`.`app_id` AND "+
			"`gluee`.`users`.`token` = `guc`.`user_id2` "+
			"WHERE `guc`.`app_id`=? AND `guc`.`user_id1`=?",
		appID,
		userToken,
	)

	return
}

// AddApplicationUserConnection will add a new connection between users or returns an error
func AddApplicationUserConnection(appID uint64, user1Token, user2Token string) (err error) {
	query := "INSERT INTO `gluee`.`user_connections` (`app_id`, `user_id1`, `user_id2`) VALUES (?, ?, ?)"
	_, err = GetMaster().
		Exec(query, appID, user1Token, user2Token)
	if err != nil {
		if config.GetConfig().Env == "dev" {
			return err
		}
		return fmt.Errorf("error while saving to database")
	}

	return
}