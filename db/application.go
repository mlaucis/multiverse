/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/tapglue/backend/entity"
)

// GetApplicationByID returns the user matching the account or an error
func GetApplicationByID(appID uint64) (application *entity.Application, err error) {
	application = &entity.Application{}

	// Execute query to get application
	err = GetSlave().
		QueryRowx("SELECT * FROM `applications` WHERE `id`=?", appID).
		StructScan(application)

	return
}

// GetAccountAllApplications returns all the users from a certain account
func GetAccountAllApplications(accountID uint64) (applications []*entity.Application, err error) {
	applications = []*entity.Application{}

	// Execute query to get applications
	err = GetSlave().
		Select(&applications, "SELECT * FROM `applications` WHERE `account_id`=?", accountID)

	return
}

// AddAccountApplication creates a user for an account and returns the created entry or an error
func AddAccountApplication(accountID uint64, application *entity.Application) (*entity.Application, error) {
	// Check if key empty
	if application.Key == "" {
		return nil, fmt.Errorf("empty application key is not allowed")
	}

	// Check if name empty
	if application.Name == "" {
		return nil, fmt.Errorf("empty application name is not allowed")
	}

	// Write to db
	query := "INSERT INTO `applications` (`account_id`, `key`, `name`) VALUES (?, ?, ?)"
	result, err := GetMaster().Exec(query, accountID, application.Key, application.Name)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	// Retrieve application
	var createdApplicationID int64
	createdApplicationID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	// Return application
	return GetApplicationByID(uint64(createdApplicationID))
}
