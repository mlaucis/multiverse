/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package entity provides all the entities needed by the app to interact
// with the API or the database
package entity

import "time"

type (
	// Commenly used structure
	Common struct {
		Image     []*Image
		Metadata  string    `json:"metadata,omitempty"`
		Enabled   bool      `json:"enabled,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	// Commonly used structure for users
	UserCommon struct {
		Username    string    `json:"user_name,omitempty"`
		Password    string    `json:"password,omitempty"`
		DisplayName string    `json:"display_name,omitempty"`
		FirstName   string    `json:"first_name,omitempty"`
		LastName    string    `json:"last_name,omitempty"`
		Email       string    `json:"email,omitempty"`
		URL         string    `json:"url,omitempty"`
		Activated   string    `json:"activated,omitempty"`
		LastLogin   time.Time `json:"last_login,omitempty"`
	}

	// Image structure
	Image struct {
		URL    string `json:"url,omitempty"`
		Type   string `json:"type,omitempty"`
		Width  string `json:"width,omitempty"`
		Heigth string `json:"height,omitempty"`
	}

	// Object structure
	Object struct {
		ID          string    `json:"id,omitempty"`
		Type        string    `json:"type,omitempty"`
		URL         string    `json:"url,omitempty"`
		DisplayName []*string `json:"display_name,omitempty"`
	}

	// Object structure
	Participant struct {
		ID          string    `json:"id,omitempty"`
		URL         string    `json:"url,omitempty"`
		DisplayName []*string `json:"display_name,omitempty"`
		Image       []*Image
	}

	// Account structure
	Account struct {
		ID          uint64 `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Enabled     bool   `json:"enabled,omitempty"`
		Common
	}

	// AccountRole structure
	AccountRole struct {
		ID          uint64 `json:"id,omitempty"`
		Permission  string `json:"permission,omitempty"`
		Description string `json:"description,omitempty"`
		Common
	}

	// AccountUser structure
	AccountUser struct {
		ID        uint64 `json:"id,omitempty"`
		AccountID uint64 `json:"account_id,omitempty"`
		Role      *AccountRole
		UserCommon
		Common
	}

	// Application structure
	Application struct {
		ID          uint64 `json:"id,omitempty"`
		AccountID   uint64 `json:"account_id,omitempty"`
		AuthToken   string `json:"auth_token,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		URL         string `json:"url,omitempty"`
		Common
	}

	// User structure
	User struct {
		ID            uint64    `json:"id",omitempty`
		ApplicationID uint64    `json:"application_id,omitempty"`
		AuthToken     string    `json:"auth_token,omitempty"`
		FacebookID    string    `json:"facebook_id,omitempty"`
		TwitterID     string    `json:"twitter_id,omitempty"`
		GoogleID      string    `json:"google_id,omitempty"`
		CustomerID    string    `json:"customer_id,omitempty"`
		GameCenterID  string    `json:"game_center_id,omitempty"`
		DeviceIDs     []*string `json:"device_ids,omitempty"`
		URL           string    `json:"url,omitempty"`
		UserCommon
		Common
	}

	// Connection structure holds the connections of the users
	Connection struct {
		ApplicationID string `json:"application_id,omitempty"`
		UserFromID    string `json:"user_from_id,omitempty"`
		UserToID      string `json:"user_to_id,omitempty"`
		Common
	}

	// Device structure
	Device struct {
		ID           string `json:"id,omitempty"`
		UUID         string `json:"uuid,omitempty"`
		IDFA         string `json:"idfa,omitempty"`
		IDFV         string `json:"idfv,omitempty"`
		GPSAdID      string `json:"gps_adid,omitempty"`
		AndroidID    string `json:"android_id,omitempty"`
		PushToken    string `json:"push_token,omitempty"`
		Mac          string `json:"mac,omitempty"`
		MacMD5       string `json:"mac_md5,omitempty"`
		MacSHA1      string `json:"mac_sha1,omitempty"`
		Platform     string `json:"platfrom,omitempty"`
		OSVersion    string `json:"os_version,omitempty"`
		Browser      string `json:"browser,omitempty"`
		Model        string `json:"model,omitempty"`
		Manufacturer string `json:"manufacturer,omitempty"`
		AppVersion   string `json:"app_version,omitempty"`
		SDKVersion   string `json:"sdk_version,omitempty"`
		Timezone     string `json:"timezone,omitempty"`
		Language     string `json:"language,omitempty"`
		Country      string `json:"country,omitempty"`
		City         string `json:"city,omitempty"`
		IP           string `json:"ip,omitempty"`
		Carrier      string `json:"carrier,omitempty"`
		Network      string `json:"network,omitempty"`
		Common
	}

	// Event structure
	Event struct {
		ID            uint64 `json:"id,omitempty"`
		ApplicationID uint64 `json:"application_id,omitempty"`
		UserID        uint64 `json:"user_id,omitempty"`
		Verb          string `json:"string,omitempty"`
		Language      string `json:"language,omitempty"`
		Prioritity    string `json:"priority,omitempty"`
		Status        string `json:"status,omitempty"`
		Location      string `json:"location,omitempty"`
		Object        *Object
		Target        *Object
		Instrument    *Object
		Attachment    *Object
		Participant   *Participant
		Common
	}
)
