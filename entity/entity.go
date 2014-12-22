/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package entity

type (

	Dates struct {
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	Application struct {
		AppID uint64 `json:"app_id"`
		Key	string `json: "key"`
		Name string `json:"name"`
		AccountID uint64 `json: "accountId"`
		Enabled bool `json: "enabled"`
		Dates
	}

	Device struct {
		DeviceID string `json:"device_id"`
		Model string `json:"model"`
		Manufacturer string `json:"manufacturer"`
		UUID string `json:"uuid"`
		IDFA string `json:"idfa"`
		AndroidID string `json:"android_id"`
		Platform string `json:"platfrom"`
		OSVersion string `json:"os_version"`
		AppVersion string `json:"app_version"`
		SDKVersion string `json:"sdk_version"`
		Timezone string `json:"timezone"`
		Language string `json:"language"`
		Country string `json:"country"`
		City string `json:"city"`
		IP string `json:"ip"`
		Carrier string `json:"carrier"`
		Network string `json:"network"`
		Enabled bool `json: "enabled"`
		Dates
	}

	Session struct {
		SessionID uint64 `"json:"session_id"`
		AppID uint64 `json:"app_id"`
		Token string  `json:"token,omitempty"`
		Nth string `json:"nth"`
		Custom string `json:"custom,omitempty"`
		Device
		Dates
	}

	User struct {
		AppID        uint64  `json:"app_id,omitempty"`
		Token        string  `json:"token,omitempty"`
		Username string `json:"username,omitempty"`
		Name string `json:"name",omitempty"`
		Password string `json:"password,omitempty"`
		Email string `json:"email,omitempty"`
		DisplayName  string  `json:"display_name,omitempty"`
		URL          string  `json:"url,omitempty"`
		ThumbnailUrl string  `json:"thumbnail_url,omitempty"`
		Custom       string  `json:"custom,omitempty"`
		LastLogin    string  `json:"last_login,omitempty"`
		Friends      []*User `json:"friends,omitempty"`
		Dates
	}

	Event struct {
		ID           uint64 `json:"eventId"`
		AppID        uint64 `json:"app_id,omitempty"`
		SessionID uint64 `"json:"session_id"`
		UserID       uint64 `json:"user_id,omitempty"`
		EventType    string `json:"event_type"`
		ThumbnailUrl string `json:"thumbnail_url,omitempty"`
		ItemID       string `json:"item_id"`
		ItemName     string `json:"item_name,omitempty"`
		ItemURL      string `json:"item_url,omitempty"`
		CreatedAt    string `json:"created_at"`
		Custom       string `json:"custom,omitempty"`
		Nth string `json:"nth"`
		User         *User  `json:"user,omitempty"` //refactor
	}

	Friend struct {
		UserID1 string `json:"user_id_1"`
		UserID2 string `json:"user_id_2"`
		Enabled bool `json: "enabled"`
		Dates
	}

	Account struct {
		AccountID uint64 `json: "accountId"`
		Name string `json: "name"`
		Enabled bool `json: "enabled"`
		Dates
	}

	AccountUser struct {
		UserId uint64 `json: "accountId"`
		AccountId uint64 `json: "accountId"`
		Name string `json: "name"`
		Password string `json: "password"`
		Email string `json: "email"`
		Enabled bool `json: "enabled"`
		Dates
	}
)
