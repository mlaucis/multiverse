/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

//Package entity provides all the entities needed by the app to interact with the database
package entity

import "time"

type (
	// Common holds common used fields
	Common struct {
		Image      []*Image  `json:"image,omitempty"`
		Metadata   string    `json:"metadata,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		ReceivedAt time.Time `json:"received_at"`
	}

	// UserCommon holds common used fields for users
	UserCommon struct {
		Username  string    `json:"user_name"`
		Password  string    `json:"password,omitempty"`
		FirstName string    `json:"first_name,omitempty"`
		LastName  string    `json:"last_name,omitempty"`
		Email     string    `json:"email,omitempty"`
		URL       string    `json:"url,omitempty"`
		Activated bool      `json:"activated,omitempty"`
		LastLogin time.Time `json:"last_login,omitempty"`
	}

	// Image structure
	Image struct {
		URL    string `json:"url"`
		Type   string `json:"type,omitempty"` // image/jpeg image/png
		Width  string `json:"width,omitempty"`
		Heigth string `json:"height,omitempty"`
	}

	// Object structure
	Object struct {
		ID          string            `json:"id"`
		Type        string            `json:"type"`
		URL         string            `json:"url,omitempty"`
		DisplayName map[string]string `json:"display_name"` // ["en"=>"article", "de"=>"artikel"]
	}

	// Participant structure
	Participant struct {
		ID    string   `json:"id"`
		URL   string   `json:"url,omitempty"`
		Image []*Image `json:"image,omitempty"`
	}

	// Account structure
	Account struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Token       string `json:"token,omitempty"`
		Common
	}

	// AccountRole structure
	AccountRole struct {
		ID          int64  `json:"id"`
		Permission  string `json:"permission"`
		Description string `json:"description"`
		Common
	}

	// AccountUser structure
	AccountUser struct {
		ID        int64        `json:"id"`
		AccountID int64        `json:"account_id"`
		Role      *AccountRole `json:"account_role,omitempty"`
		UserCommon
		Common
	}

	// Application structure
	Application struct {
		ID          int64  `json:"id"`
		AccountID   int64  `json:"account_id"`
		AuthToken   string `json:"auth_token"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Common
	}

	// User structure
	User struct {
		ID            int64     `json:"id"`
		AccountID     int64     `json:"account_id"`
		ApplicationID int64     `json:"application_id"`
		AuthToken     string    `json:"auth_token"`
		FacebookID    string    `json:"facebook_id,omitempty"`
		TwitterID     string    `json:"twitter_id,omitempty"`
		GoogleID      string    `json:"google_id,omitempty"`
		CustomerID    string    `json:"customer_id,omitempty"`
		GameCenterID  string    `json:"game_center_id,omitempty"`
		DeviceIDs     []*string `json:"device_ids,omitempty"`
		Events        []*Event  `json:"events,omitempty"`
		Connections   []*User   `json:"connections,omitempty"`
		UserCommon
		Common
	}

	// Connection structure holds the connections of the users
	Connection struct {
		AccountID     int64 `json:"account_id"`
		ApplicationID int64 `json:"application_id"`
		UserFromID    int64 `json:"user_from_id"`
		UserToID      int64 `json:"user_to_id"`
		Common
	}

	// Device structure
	Device struct {
		ID           string `json:"id"`
		UUID         string `json:"uuid,omitempty"`
		IDFA         string `json:"idfa,omitempty"`
		IDFV         string `json:"idfv,omitempty"`
		GPSAdID      string `json:"gps_adid,omitempty"`
		AndroidID    string `json:"android_id,omitempty"`
		PushToken    string `json:"push_token,omitempty"`
		MAC          string `json:"mac,omitempty"`
		MACMD5       string `json:"mac_md5,omitempty"`
		MACSHA1      string `json:"mac_sha1,omitempty"`
		Platform     string `json:"platfrom"`
		OSVersion    string `json:"os_version"`
		Browser      string `json:"browser,omitempty"`
		Model        string `json:"model"`
		Manufacturer string `json:"manufacturer"`
		AppVersion   string `json:"app_version"`
		SDKVersion   string `json:"sdk_version"`
		Timezone     string `json:"timezone"`
		Language     string `json:"language"`
		Country      string `json:"country,omitempty"`
		City         string `json:"city,omitempty"`
		IP           string `json:"ip,omitempty"`
		Carrier      string `json:"carrier,omitempty"`
		Network      string `json:"network,omitempty"`
		Common
	}

	// Event structure
	Event struct {
		ID            int64          `json:"id"`
		AccountID     int64          `json:"account_id"`
		ApplicationID int64          `json:"application_id"`
		UserID        int64          `json:"user_id"`
		Verb          string         `json:"verb"`
		Language      string         `json:"language"`
		Prioritity    string         `json:"priority"`
		Location      string         `json:"location,omitempty"`
		Object        *Object        `json:"object"`
		Target        *Object        `json:"target,omitempty"`
		Instrument    *Object        `json:"instrument,omitempty"`
		Participant   []*Participant `json:"participant,omitempty"`
		Common
	}
)
