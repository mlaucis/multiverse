/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package entity provides all the entities needed by the app to interact
// with the API or the database
package entity

import "time"

type (

	// Account structure
	Account struct {
		ID         uint64         `json:"id"`
		Name       string         `json:"name"`
		AcountUser []*AccountUser `json:"account_user,omitempty" db:"-"`
		Enabled    bool           `json:"enabled"`
		CreatedAt  time.Time      `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	}

	// AccountUser structure
	AccountUser struct {
		ID        uint64    `json:"id"`
		AccountID uint64    `json:"account_id" db:"account_id"`
		Name      string    `json:"name"`
		Password  string    `json:"password"`
		Email     string    `json:"email"`
		Enabled   bool      `json:"enabled"`
		LastLogin time.Time `json:"last_login,omitempty" db:"last_login"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	// Application structure
	Application struct {
		ID        uint64    `json:"id"`
		Key       string    `json:"key"`
		Name      string    `json:"name"`
		AccountID uint64    `json:"account_id" db:"account_id"`
		Enabled   bool      `json:"enabled"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	// User structure
	User struct {
		AppID        uint64    `json:"app_id,omitempty" db:"application_id"`
		Token        string    `json:"token,omitempty"`
		Username     string    `json:"username,omitempty"`
		Name         string    `json:"name",omitempty"`
		Password     string    `json:"password,omitempty"`
		Email        string    `json:"email,omitempty"`
		URL          string    `json:"url,omitempty" db:"url"`
		ThumbnailURL string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
		Provider     string    `json:"provider,omitempty"`
		Custom       string    `json:"custom,omitempty"`
		Connections  []*User   `json:"connections,omitempty" db:"-"`
		Events       []*Event  `json:"events,omitempty" db:"-"`
		LastLogin    time.Time `json:"last_login,omitempty" db:"last_login"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

	// UserConnection structure holds the connections of the users between each-other
	UserConnection struct {
		AppID      string    `json:"app_id" db:"application_id"`
		UserToken1 string    `json:"user_id1" db:"user_id1"`
		UserToken2 string    `json:"user_id2" db:"user_id2"`
		Enabled    bool      `json:"enabled"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	// Device structure
	Device struct {
		GID          string `json:"gid" db:"gid"`
		Model        string `json:"model"`
		Manufacturer string `json:"manufacturer"`
		UUID         string `json:"uuid" db:"uuid"`
		IDFA         string `json:"idfa" db:"idfa"`
		IDFV         string `json:"idfv" db:"idfv"`
		Mac          string `json:"mac" db:"mac"`
		MacMD5       string `json:"mac_md5" db:"mac_md5"`
		MacSHA1      string `json:"mac_sha1" db:"mac_sha1"`
		AndroidID    string `json:"android_id" db:"android_id"`
		GPSAdID      string `json:"gps_adid" db:"gps_adid"`
		Platform     string `json:"platfrom"`
		OSVersion    string `json:"os_version" db:"os_version"`
		Browser      string `json:"browser" db:"browser"`
		AppVersion   string `json:"app_version" db:"app_version"`
		SDKVersion   string `json:"sdk_version" db:"sdk_version"`
		Timezone     string `json:"timezone"`
		Language     string `json:"language"`
		Country      string `json:"country"`
		City         string `json:"city"`
		IP           string `json:"ip" db:"ip"`
		Carrier      string `json:"carrier"`
		Network      string `json:"network"`
	}

	// Session structure
	Session struct {
		ID        uint64 `json:"id"`
		AppID     uint64 `json:"app_id" db:"application_id"`
		UserToken string `json:"user_token" db:"user_token"`
		Nth       uint64 `json:"nth"`
		Custom    string `json:"custom,omitempty"`
		Device
		User      *User     `json:"user,omitempty" db:"-"`
		Events    []*Event  `json:"events,omitempty" db:"-"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	// Item structure
	Item struct {
		ID           string `json:"item_id" db:"item_id"`
		Name         string `json:"item_name,omitempty" db:"item_name"`
		URL          string `json:"item_url,omitempty" db:"item_url"`
		ThumbnailURL string `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	}

	// Event structure
	Event struct {
		ID        uint64 `json:"id"`
		AppID     uint64 `json:"app_id,omitempty" db:"application_id"`
		SessionID uint64 `json:"session_id" db:"session_id"`
		UserToken string `json:"user_token,omitempty" db:"user_token"`
		Title     string `json:"title",omitempty" db:"title"`
		Type      string `json:"type"`
		Item
		Custom    string    `json:"custom,omitempty"`
		Nth       uint64    `json:"nth"`
		User      *User     `json:"user,omitempty" db:"-"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)
