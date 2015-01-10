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

	err = GetSlave().
		QueryRowx("SELECT * FROM `applications` WHERE `id`=?", appID).
		StructScan(application)

	return
}

// GetAccountAllApplications returns all the users from a certain account
func GetAccountAllApplications(accountID uint64) (applications []*entity.Application, err error) {
	applications = []*entity.Application{}

	err = GetSlave().
		Select(&applications, "SELECT * FROM `applications` WHERE `account_id`=?", accountID)

	return
}

// AddAccountApplication creates a user for an account and returns the created entry or an error
func AddAccountApplication(accountID uint64, application *entity.Application) (*entity.Application, error) {
	if application.Key == "" {
		return nil, fmt.Errorf("empty application key is not allowed")
	}

	if application.Name == "" {
		return nil, fmt.Errorf("empty application name is not allowed")
	}

	query := "INSERT INTO `applications` (`account_id`, `key`, `name`) VALUES (?, ?, ?)"
	result, err := GetMaster().Exec(query, accountID, application.Key, application.Name)
	if err != nil {
		return nil, fmt.Errorf("error while saving to database")
	}

	var createdApplicationID int64
	createdApplicationID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	return GetApplicationByID(uint64(createdApplicationID))
}
