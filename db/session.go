/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package db

import (
	"fmt"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/entity"
)

// GetSessionByID returns the session matching the id or an error
func GetSessionByID(sessionID uint64) (session *entity.Session, err error) {
	session = &entity.Session{}

	err = GetSlave().
		QueryRowx("SELECT * FROM `sessions` WHERE `id`=?", sessionID).
		StructScan(session)

	return
}

// GetAllUserSessions returns all the sessions for a certain user
func GetAllUserSessions(appID uint64, userToken string) (userSessions []*entity.Session, err error) {
	userSessions = []*entity.Session{}

	err = GetSlave().
		Select(&userSessions, "SELECT * FROM `sessions` WHERE `application_id`=? AND `user_token`=?", appID, userToken)

	return
}

// AddUserSession creates a new session for an user and returns the created entry or an error
func AddUserSession(session *entity.Session) (*entity.Session, error) {
	query := "INSERT INTO `sessions` (`application_id`, `user_token`, `nth`, `custom`, " +
		"`language`, `country`, `network`, `uuid`, `platform`, `sdk_version`, `timezone`, `city`, `gid`, " +
		"`idfa`, `idfv`, `mac`, `mac_md5`, `mac_sha1`, `gps_adid`, `app_version`, `carrier`, `model`, `manufacturer`, `android_id`, `os_version`, `ip`, `browser`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := GetMaster().
		Exec(query,
		session.AppID,
		session.UserToken,
		session.Nth,
		session.Custom,
		session.Device.Language,
		session.Device.Country,
		session.Device.Network,
		session.Device.UUID,
		session.Device.Platform,
		session.Device.SDKVersion,
		session.Device.Timezone,
		session.Device.City,
		session.Device.GID,
		session.Device.IDFA,
		session.Device.IDFV,
		session.Device.Mac,
		session.Device.MacMD5,
		session.Device.MacSHA1,
		session.Device.GPSAdID,
		session.Device.AppVersion,
		session.Device.Carrier,
		session.Device.Model,
		session.Device.Manufacturer,
		session.Device.AndroidID,
		session.Device.OSVersion,
		session.Device.IP,
		session.Device.Browser,
	)

	if err != nil {
		if config.Conf().Env() == "dev" {
			return nil, err
		}
		return nil, fmt.Errorf("error while saving to database")
	}

	var createdSessionID int64
	createdSessionID, err = result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error while processing the request")
	}

	return GetSessionByID(uint64(createdSessionID))
}
