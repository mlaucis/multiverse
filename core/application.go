/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

// ReadApplication returns the application matching the ID or an error
func ReadApplication(accountID, applicationID int64) (application *entity.Application, err error) {
	result, err := storageEngine.Get(storageClient.AccountAppKey(accountID, applicationID)).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &application); err != nil {
		return nil, err
	}

	return
}

// ReadApplicationList returns all applications from a certain account
func ReadApplicationList(accountID int64) (applications []*entity.Application, err error) {
	key := storageClient.AccountAppsKey(accountID)

	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		err := errors.New("There are no apps for this account")
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	application := &entity.Application{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), application); err != nil {
			return nil, err
		}
		applications = append(applications, application)
		application = &entity.Application{}
	}

	return
}

// WriteApplication adds a application to the database and returns the created applicaton user or an error
func WriteApplication(application *entity.Application, retrieve bool) (app *entity.Application, err error) {
	if application.ID, err = storageClient.GenerateApplicationID(application.AccountID); err != nil {
		return nil, err
	}

	val, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}

	key := storageClient.AccountAppKey(application.AccountID, application.ID)

	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("application already exists")
	}
	if err != nil {
		return nil, err
	}

	// Generate list key
	listKey := storageClient.AccountAppsKey(application.AccountID)

	// Write list
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return application, nil
	}

	// Return resource
	return ReadApplication(application.AccountID, application.ID)
}
