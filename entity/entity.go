/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package entity

type (
	Dates struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	Device struct {
		DeviceID     string `json:"device_id"`
		Model        string `json:"model"`
		Manufacturer string `json:"manufacturer"`
		UUID         string `json:"uuid"`
		IDFA         string `json:"idfa"`
		AndroidID    string `json:"android_id"`
		Platform     string `json:"platfrom"`
		OSVersion    string `json:"os_version"`
		AppVersion   string `json:"app_version"`
		SDKVersion   string `json:"sdk_version"`
		Timezone     string `json:"timezone"`
		Language     string `json:"language"`
		Country      string `json:"country"`
		City         string `json:"city"`
		IP           string `json:"ip"`
		Carrier      string `json:"carrier"`
		Network      string `json:"network"`
		Enabled      bool   `json: "enabled"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}

	Account struct {
		AccountID  uint64         `json: "account_id"`
		Name       string         `json: "name"`
		AcountUser []*AccountUser `json:"account_user,omitempty"`
		Enabled    bool           `json: "enabled"`
		CreatedAt  string         `json:"created_at"`
		UpdatedAt  string         `json:"updated_at"`
	}

	AccountUser struct {
		UserID    string `json: "user_id"`
		AccountID uint64 `json: "account_id"`
		Name      string `json: "name"`
		Password  string `json: "password"`
		Email     string `json: "email"`
		Enabled   bool   `json: "enabled"`
		LastLogin string `json:"last_login,omitempty"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	Application struct {
		AppID     uint64 `json:"app_id"`
		Key       string `json: "key"`
		Name      string `json:"name"`
		AccountID uint64 `json: "account_id"`
		Enabled   bool   `json: "enabled"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	User struct {
		AppID        uint64  `json:"app_id,omitempty"`
		Token        string  `json:"token,omitempty"`
		Username     string  `json:"username,omitempty"`
		Name         string  `json:"name",omitempty"`
		Password     string  `json:"password,omitempty"`
		Email        string  `json:"email,omitempty"`
		URL          string  `json:"url,omitempty"`
		ThumbnailUrl string  `json:"thumbnail_url,omitempty"`
		Custom       string  `json:"custom,omitempty"`
		Connections  []*User `json:"connections,omitempty"`
		LastLogin    string  `json:"last_login,omitempty"`
		CreatedAt    string  `json:"created_at"`
		UpdatedAt    string  `json:"updated_at"`
	}

	UserConnection struct {
		UserID1   string `json:"user_id_1"`
		UserID2   string `json:"user_id_2"`
		Enabled   bool   `json: "enabled"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	Session struct {
		SessionID uint64 `"json:"session_id"`
		AppID     uint64 `json:"app_id"`
		Token     string `json:"token,omitempty"`
		Nth       string `json:"nth"`
		Custom    string `json:"custom,omitempty"`
		Device
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	Event struct {
		EventID      uint64 `json:"event_id"`
		AppID        uint64 `json:"app_id,omitempty"`
		SessionID    uint64 `"json:"session_id"`
		UserID       uint64 `json:"user_id,omitempty"`
		EventType    string `json:"event_type"`
		ThumbnailUrl string `json:"thumbnail_url,omitempty"`
		ItemID       string `json:"item_id"`
		ItemName     string `json:"item_name,omitempty"`
		ItemURL      string `json:"item_url,omitempty"`
		CreatedAt    string `json:"created_at"`
		Custom       string `json:"custom,omitempty"`
		Nth          string `json:"nth"`
		User         *User  `json:"user,omitempty"` //refactor
	}
)
