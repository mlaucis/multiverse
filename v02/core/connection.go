/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Connection interface
	Connection interface {
		Create(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)
		Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err tgerrors.TGError)
		Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)
		Delete(accountID, applicationID, userFromID, userToID int64) (err tgerrors.TGError)
		List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)
		FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)
		Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)
		WriteEventsToList(connection *entity.Connection) (err tgerrors.TGError)
		DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err tgerrors.TGError)
		SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string) ([]*entity.ApplicationUser, tgerrors.TGError)
		AutoConnectSocialFriends(user *entity.ApplicationUser, ourStoredUsersIDs []interface{}) (users []*entity.ApplicationUser, err tgerrors.TGError)
	}
)
