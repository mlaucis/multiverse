/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/redis"
)

// Defining keys
const (
	ApplicationKey  string = "account_%d_application_%d"
	ApplicationsKey string = "account_%d_applications"
)

// generateApplicationID generates a new application ID
func generateApplicationID(accountID int64) (int64, error) {
	incr := redis.Client().Incr(fmt.Sprintf("ids_account_%d_application", accountID))
	return incr.Result()
}

// ReadApplication returns the application matching the ID or an error
func ReadApplication(accountID int64, applicationID int64) (application *entity.Application, err error) {
	// Generate resource key
	key := fmt.Sprintf(ApplicationKey, accountID, applicationID)

	// Read from db
	result, err := redis.Client().Get(key).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &application); err != nil {
		return nil, err
	}

	return
}

// ReadApplicationList returns all applications from a certain account
func ReadApplicationList(accountID int64) (applications []*entity.Application, err error) {
	// Generate resource key
	key := fmt.Sprintf(ApplicationsKey, accountID)

	// Read from db
	result, err := redis.Client().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no apps for this account")
		return nil, err
	}

	// Read from db
	resultList, err := redis.Client().MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
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
	// Generate id
	if application.ID, err = generateApplicationID(application.AccountID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(application)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := fmt.Sprintf(ApplicationKey, application.AccountID, application.ID)

	// Write resource
	if err = redis.Client().Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := fmt.Sprintf(ApplicationsKey, application.AccountID)

	// Write list
	if err = redis.Client().LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return application, nil
	}

	// Return resource
	return ReadApplication(application.AccountID, application.ID)
}
