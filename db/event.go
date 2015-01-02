/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/entity"
)

// GetEventByID returns the event matching the id or an error
func GetEventByID(eventID uint64) (event *entity.Event, err error) {
	event = &entity.Event{}

	err = GetSlave().
		QueryRowx("SELECT * FROM `tapglue`.`events` WHERE `id`=?", eventID).
		StructScan(event)

	return
}

// GetAllUserAppEvents returns all the events of a user for a certain app
func GetAllUserAppEvents(appID uint64, userToken string) (user *entity.User, err error) {
	if user, err = GetApplicationUserByToken(appID, userToken); err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return
	}

	user.Events = []*entity.Event{}

	err = GetSlave().
		Select(&user.Events, "SELECT * FROM `tapglue`.`events` WHERE `application_id`=? AND `user_token`=?", appID, userToken)

	return
}

// GetSessionEvents returns all the events from a session of a user within an app
func GetSessionEvents(appID, sessionID uint64, userToken string) (session *entity.Session, err error) {
	if session, err = GetSessionByID(sessionID); err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return
	}

	if session.AppID != appID || session.UserToken != userToken {
		return nil, fmt.Errorf("invalid session retrieved")
	}

	session.User = &entity.User{}

	if session.User, err = GetApplicationUserByToken(appID, userToken); err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return
	}

	session.Events = []*entity.Event{}

	err = GetSlave().
		Select(&session.Events, "SELECT * FROM `tapglue`.`events` WHERE `application_id`=? AND `session_id`=? AND `user_token`=?", appID, sessionID, userToken)

	return
}

// GetUserConnectionsEvents returns the events of all the users connected to the specified user
func GetUserConnectionsEvents(appID uint64, userToken string) (events []*entity.Event, err error) {
	events = []*entity.Event{}

	query := "SELECT `tapglue`.`events`.* " +
		"FROM `tapglue`.`events` " +
		"JOIN `tapglue`.`user_connections` as `guc` " +
		"ON `tapglue`.`events`.`application_id` = `guc`.`application_id`" +
		"AND `tapglue`.`events`.`user_token` = `guc`.`user_id2` " +
		"WHERE `guc`.`application_id`=? AND `guc`.`user_id1`=?"

	err = GetSlave().
		Select(&events, query, appID, userToken)

	return
}

// AddSessionEvent creates a new session for an user and returns the created entry or an error
func AddSessionEvent(event *entity.Event) (*entity.Event, error) {
	query := "INSERT INTO `tapglue`.`events` (`application_id`, `session_id`, `user_token`, `type`, " +
		"`item_id`, `item_name`, `item_url`, `thumbnail_url`, `custom`, `nth`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := GetMaster().
		Exec(query,
		event.AppID,
		event.SessionID,
		event.UserToken,
		event.Type,
		event.Item.ID,
		event.Item.Name,
		event.Item.URL,
		event.Item.ThumbnailURL,
		event.Custom,
		event.Nth,
	)

	if err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return nil, fmt.Errorf("error while saving to database")
	}

	var createdEventID int64
	createdEventID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	return GetEventByID(uint64(createdEventID))
}
