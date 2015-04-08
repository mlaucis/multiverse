/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/entity"
)

// ReadApplication returns the application matching the ID or an error
func ReadApplication(accountID, applicationID int64) (*entity.Application, *tgerrors.TGError) {
	result, er := redisEngine.Get(storageClient.Application(accountID, applicationID)).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application (1)", er.Error())
	}

	application := &entity.Application{}
	if er := json.Unmarshal([]byte(result), application); er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application (2)", er.Error())
	}

	return application, nil
}

// UpdateApplication updates an application in the database and returns the created applicaton user or an error
func UpdateApplication(existingApplication, updatedApplication entity.Application, retrieve bool) (*entity.Application, *tgerrors.TGError) {
	updatedApplication.UpdatedAt = time.Now()

	val, er := json.Marshal(updatedApplication)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application (1)\n"+er.Error(), er.Error())
	}

	key := storageClient.Application(updatedApplication.AccountID, updatedApplication.ID)
	exist, er := redisEngine.Exists(key).Result()
	if !exist {
		return nil, tgerrors.NewNotFoundError("failed to update the application (2)", "app not found")
	}
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application (3)", er.Error())
	}

	if er = redisEngine.Set(key, string(val)).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application (4)", er.Error())
	}

	if !updatedApplication.Enabled {
		listKey := storageClient.Applications(updatedApplication.AccountID)
		if er = redisEngine.LRem(listKey, 0, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application (5)", er.Error())
		}
	}

	if !retrieve {
		return &updatedApplication, nil
	}

	return ReadApplication(updatedApplication.AccountID, updatedApplication.ID)
}

// DeleteApplication deletes the application matching the IDs or an error
func DeleteApplication(accountID, applicationID int64) *tgerrors.TGError {
	// TODO: Disable application users?
	// TODO: User connections?
	// TODO: Application lists?
	// TODO: Application events?

	key := storageClient.Application(accountID, applicationID)
	result, er := redisEngine.Del(key).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to delete the application (1)", er.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to delete the application (2)", "app not found")
	}

	listKey := storageClient.Applications(accountID)
	if er := redisEngine.LRem(listKey, 0, key).Err(); er != nil {
		return tgerrors.NewInternalError("failed to delete the application (3)", er.Error())
	}

	return nil
}

// ReadApplicationList returns all applications from a certain account
func ReadApplicationList(accountID int64) ([]*entity.Application, *tgerrors.TGError) {
	key := storageClient.Applications(accountID)

	result, er := redisEngine.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the applications list (1)", er.Error())
	}

	applications := []*entity.Application{}
	if len(result) == 0 {
		return applications, nil
	}

	resultList, er := redisEngine.MGet(result...).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the applications list (2)", er.Error())
	}

	application := &entity.Application{}
	for _, result := range resultList {
		if er = json.Unmarshal([]byte(result.(string)), application); er != nil {
			return nil, tgerrors.NewInternalError("failed to read the applications list (3)", er.Error())
		}
		applications = append(applications, application)
		application = &entity.Application{}
	}

	return applications, nil
}

// WriteApplication adds an application to the database and returns the created applicaton user or an error
func WriteApplication(application *entity.Application, retrieve bool) (*entity.Application, *tgerrors.TGError) {
	var er error
	if application.ID, er = storageClient.GenerateApplicationID(application.AccountID); er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (1)", er.Error())
	}

	application.Enabled = true
	application.CreatedAt = time.Now()
	application.UpdatedAt = application.CreatedAt

	if application.AuthToken, er = storageClient.GenerateApplicationSecretKey(application); er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (2)", er.Error())
	}

	val, er := json.Marshal(application)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (3)", er.Error())
	}

	key := storageClient.Application(application.AccountID, application.ID)

	exist, er := redisEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, tgerrors.NewInternalError("failed to create the application (3)", "duplicate app")
	}
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (4)", er.Error())
	}

	listKey := storageClient.Applications(application.AccountID)

	if er = redisEngine.LPush(listKey, key).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (5)", er.Error())
	}

	// Store the token details in redis
	_, er = redisEngine.HMSet(
		"tokens:"+utils.Base64Encode(application.AuthToken),
		"acc", strconv.FormatInt(application.AccountID, 10),
		"app", strconv.FormatInt(application.ID, 10),
	).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to create the application (6)", er.Error())
	}

	if !retrieve {
		return application, nil
	}

	return ReadApplication(application.AccountID, application.ID)
}
