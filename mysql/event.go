/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package mysql

import (
	"fmt"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/entity"
)

// GetEventByID returns the event matching the id or an error
func GetEventByID(eventID uint64) (event *entity.Event, err error) {
	event = &entity.Event{}

	// Exectute query to get event
	err = GetSlave().
		QueryRowx("SELECT * FROM `events` WHERE `id`=?", eventID).
		StructScan(event)

	return
}

// GetAllUserAppEvents returns all the events of a user for a certain app
func GetAllUserAppEvents(appID uint64, userID string) (user *entity.User, err error) {
	if user, err = GetApplicationUserByToken(appID, userID); err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return
	}

	user.Events = []*entity.Event{}

	// Execture query to get events
	err = GetSlave().
		Select(&user.Events, "SELECT * FROM `events` WHERE `application_id`=? AND `user_token`=?", appID, userID)

	return
}

// GetUserConnectionsEvents returns the events of all the users connected to the specified user
func GetUserConnectionsEvents(appID uint64, userToken string) (events []*entity.Event, err error) {
	events = []*entity.Event{}

	// Execute query to get connection events
	query := "SELECT `events`.* " +
		"FROM `events` " +
		"JOIN `user_connections` as `guc` " +
		"ON `events`.`application_id` = `guc`.`application_id`" +
		"AND `events`.`user_token` = `guc`.`user_id2` " +
		"WHERE `guc`.`application_id`=? AND `guc`.`user_id1`=?"

	err = GetSlave().
		Select(&events, query, appID, userToken)

	return
}

// AddSessionEvent creates a new session for an user and returns the created entry or an error
func AddSessionEvent(event *entity.Event) (*entity.Event, error) {
	// Execute query to write an event
	query := "INSERT INTO `events` (`application_id`, `session_id`, `user_token`, `title`, `type`, " +
		"`item_id`, `item_name`, `item_url`, `thumbnail_url`, `custom`, `nth`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := GetMaster().
		Exec(query,
		event.ApplicationID,
		event.UserID,
		event.Object.DisplayName["en"],
		event.Verb,
		event.Object.ID,
		event.Object.URL,
		event.Image[0].URL,
		event.Metadata,
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

	// Return event
	return GetEventByID(uint64(createdEventID))
}
