/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package entity provides all the entities needed by the app to interact
// with the API or the database
package entity

import "time"

type (
	// Device structure
	Device struct {
		GID          string    `json:"gid" db:"gid"`
		Model        string    `json:"model"`
		Manufacturer string    `json:"manufacturer"`
		UUID         string    `json:"uuid" db:"uuid"`
		IDFA         string    `json:"idfa" db:"idfa"`
		AndroidID    string    `json:"android_id" db:"android_id"`
		Platform     string    `json:"platfrom"`
		OSVersion    string    `json:"os_version" db:"os_version"`
		AppVersion   string    `json:"app_version" db:"app_version"`
		SDKVersion   string    `json:"sdk_version" db:"sdk_version"`
		Timezone     string    `json:"timezone"`
		Language     string    `json:"language"`
		Country      string    `json:"country"`
		City         string    `json:"city"`
		IP           string    `json:"ip" db:"ip"`
		Carrier      string    `json:"carrier"`
		Network      string    `json:"network"`
		Enabled      bool      `json:"enabled"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

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
		ID        string    `json:"id"`
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
		AppID        uint64    `json:"app_id,omitempty" db:"app_id"`
		Token        string    `json:"token,omitempty"`
		Username     string    `json:"username,omitempty"`
		Name         string    `json:"name",omitempty"`
		Password     string    `json:"password,omitempty"`
		Email        string    `json:"email,omitempty"`
		URL          string    `json:"url,omitempty"`
		ThumbnailURL string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
		Custom       string    `json:"custom,omitempty"`
		Connections  []*User   `json:"connections,omitempty" db:"-"`
		LastLogin    time.Time `json:"last_login,omitempty" db:"last_login"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

	// UserConnection structure holds the connections of the users between each-other
	UserConnection struct {
		AppID      string    `json:"app_id" db:"app_id"`
		UserToken1 string    `json:"user_token1" db:"user_token1"`
		UserToken2 string    `json:"user_token2" db:"user_token2"`
		Enabled    bool      `json:"enabled"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	// Session structure
	Session struct {
		ID        uint64    `json:"id"`
		AppID     uint64    `json:"app_id" db:"app_id"`
		Token     string    `json:"token"`
		Nth       uint64    `json:"nth"`
		Custom    string    `json:"custom,omitempty"`
		Device    *Device   `json:"device"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	}

	// Item structure
	Item struct {
		ID   string `json:"id"`
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty" db:"url"`
	}

	// Event structure
	Event struct {
		ID           uint64    `json:"id"`
		AppID        uint64    `json:"app_id,omitempty" db:"app_id"`
		SessionID    uint64    `json:"session_id" db:"session_id"`
		UserToken    uint64    `json:"user_token,omitempty" db:"user_token"`
		Type         string    `json:"type"`
		ThumbnailURL string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
		Item         *Item     `json:"item"`
		Custom       string    `json:"custom,omitempty"`
		Nth          uint64    `json:"nth"`
		User         *User     `json:"user,omitempty" db:"-"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
	}
)
