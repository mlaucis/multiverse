/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

//Package entity provides all the entities needed by the app to interact with the database
package entity

import (
	"time"
)

type (
	// Common holds common used fields
	Common struct {
		Metadata  map[string]string `json:"metadata,omitempty"`
		Image     []*Image          `json:"image,omitempty"`
		CreatedAt *time.Time        `json:"created_at,omitempty"`
		UpdatedAt *time.Time        `json:"updated_at,omitempty"`
		Enabled   bool              `json:"enabled"`
	}

	// UserCommon holds common used fields for users
	UserCommon struct {
		Username         string     `json:"user_name"`
		OriginalPassword string     `json:"-"`
		Password         string     `json:"password,omitempty"`
		FirstName        string     `json:"first_name,omitempty"`
		LastName         string     `json:"last_name,omitempty"`
		Email            string     `json:"email,omitempty"`
		URL              string     `json:"url,omitempty"`
		LastLogin        *time.Time `json:"last_login,omitempty"`
		Activated        bool       `json:"activated,omitempty"`
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
		ID           int64          `json:"-"`
		PublicID     string         `json:"id"`
		Name         string         `json:"name"`
		Description  string         `json:"description"`
		AuthToken    string         `json:"token"`
		Users        []*AccountUser `json:"-"`
		Applications []*Application `json:"-"`
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
		ID              int64        `json:"-"`
		AccountID       int64        `json:"-"`
		PublicID        string       `json:"id"`
		PublicAccountID string       `json:"account_id"`
		Role            *AccountRole `json:"account_role,omitempty"`
		SessionToken    string       `json:"-"`
		UserCommon
		Common
	}

	// Application structure
	Application struct {
		ID              int64              `json:"-"`
		AccountID       int64              `json:"-"`
		PublicID        string             `json:"id"`
		PublicAccountID string             `json:"account_id"`
		AuthToken       string             `json:"token"`
		Name            string             `json:"name"`
		Description     string             `json:"description"`
		URL             string             `json:"url"`
		Users           []*ApplicationUser `json:"-"`
		Common
	}

	// ApplicationUser structure
	ApplicationUser struct {
		ID                   string              `json:"id"`
		CustomID             string              `json:"custom_id,omitempty"`
		SessionToken         string              `json:"-"`
		SocialIDs            map[string]string   `json:"social_ids,omitempty"`
		SocialConnectionsIDs map[string][]string `json:"social_connections_ids,omitempty"`
		SocialConnectionType string              `json:"connection_type,omitempty"`
		DeviceIDs            []string            `json:"device_ids,omitempty"`
		Events               []*Event            `json:"events,omitempty"`
		Connections          []*ApplicationUser  `json:"connections,omitempty"`
		LastRead             *time.Time          `json:"-"`
		UserCommon
		Common
	}

	// Connection structure holds the connections of the users
	Connection struct {
		UserFromID  string     `json:"user_from_id"`
		UserToID    string     `json:"user_to_id"`
		Type        string     `json:"type"`
		ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
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
		ID                 string         `json:"id"`
		UserID             string         `json:"user_id"`
		Verb               string         `json:"verb"`
		Language           string         `json:"language,omitempty"`
		Priority           string         `json:"priority,omitempty"`
		Location           string         `json:"location,omitempty"`
		Latitude           float64        `json:"latitude,omitempty"`
		Longitude          float64        `json:"longitude,omitempty"`
		DistanceFromTarget float64        `json:"-"`
		Visibility         uint8          `json:"visibility"`
		Object             *Object        `json:"object"`
		Target             *Object        `json:"target,omitempty"`
		Instrument         *Object        `json:"instrument,omitempty"`
		Participant        []*Participant `json:"participant,omitempty"`
		Common
	}

	// LoginPayload defines how the login payload should look like
	LoginPayload struct {
		Email     string `json:"email,omitempty"`
		Username  string `json:"user_name,omitempty"`
		EmailName string `json:"username,omitempty"`
		Password  string `json:"password"`
	}

	// SortableEventsByDistance provides the struct needed for sorting the elements by distance from target
	SortableEventsByDistance []*Event
)

const (
	// EventPrivate flags that the event is private
	EventPrivate = 10

	// EventConnections flags that the event is shared with the connetions of the user
	EventConnections = 20

	// EventPublic flags that the event is public
	EventPublic = 30
)

func (e SortableEventsByDistance) Len() int      { return len(e) }
func (e SortableEventsByDistance) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e SortableEventsByDistance) Less(i, j int) bool {
	return e[i].DistanceFromTarget < e[j].DistanceFromTarget
}
