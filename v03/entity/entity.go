//Package entity provides all the entities needed by the app to interact with the database
package entity

import (
	"time"
)

type (
	// AccAppIDs holds the account and application IDs
	AccAppIDs struct {
		AccountID           int64  `json:"acc_id,omitempty"`
		ApplicationID       int64  `json:"app_id,omitempty"`
		PublicAccountID     string `json:"pub_acc_id,omitempty"`
		PublicApplicationID string `json:"pub_app_id,omitempty"`
		CurrentUserID       uint64 `json:"current_user_id,omitempty"`
	}

	// Common holds common used fields
	Common struct {
		Metadata  interface{}       `json:"metadata,omitempty"`
		Images    map[string]*Image `json:"images,omitempty"`
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
		Deleted          *bool      `json:"deleted,omitempty"`
		SessionToken     string     `json:"session_token,omitempty"`
	}

	// Image structure
	Image struct {
		URL    string `json:"url"`
		Type   string `json:"type,omitempty"` // image/jpeg image/png
		Width  int    `json:"width,omitempty"`
		Heigth int    `json:"height,omitempty"`
	}

	// Object structure
	Object struct {
		ID           string            `json:"id"`
		Type         string            `json:"type"`
		URL          string            `json:"url,omitempty"`
		DisplayNames map[string]string `json:"display_names"` // ["en"=>"article", "de"=>"artikel"]
	}

	// Participant structure
	Participant struct {
		ID     string            `json:"id"`
		URL    string            `json:"url,omitempty"`
		Images map[string]*Image `json:"images,omitempty"`
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

	// AccountUser structure
	AccountUser struct {
		ID              int64  `json:"-"`
		AccountID       int64  `json:"-"`
		PublicID        string `json:"id"`
		PublicAccountID string `json:"account_id"`
		SessionToken    string `json:"-"`
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
		ID                   uint64              `json:"id"`
		CustomID             string              `json:"custom_id,omitempty"`
		SocialIDs            map[string]string   `json:"social_ids,omitempty"`
		SocialConnectionsIDs map[string][]string `json:"social_connections_ids,omitempty"`
		SocialConnectionType string              `json:"connection_type,omitempty"`
		DeviceIDs            []string            `json:"device_ids,omitempty"`
		Events               []*Event            `json:"events,omitempty"`
		Connections          []*ApplicationUser  `json:"connections,omitempty"`
		LastRead             *time.Time          `json:"-"`
		Relation
		UserCommon
		Common
	}

	// ApplicationUserWithIDs holds the application user structure with the added account and application ids
	ApplicationUserWithIDs struct {
		AccAppIDs
		ApplicationUser
	}

	// Connection structure holds the connections of the users
	Connection struct {
		UserFromID  uint64     `json:"user_from_id"`
		UserToID    uint64     `json:"user_to_id"`
		Type        string     `json:"type"`
		ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
		Common
	}

	// ConnectionWithIDs holds the connection structure with the added account and application ids
	ConnectionWithIDs struct {
		AccAppIDs
		Connection
	}

	// Relation holds the relation between two users
	Relation struct {
		IsFriends   *bool `json:"is_friends,omitempty"`
		IsFollower  *bool `json:"is_follower,omitempty"`
		IsFollowing *bool `json:"is_following,omitempty"`
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
		ID                 uint64         `json:"id"`
		UserID             uint64         `json:"user_id"`
		Type               string         `json:"type"`
		Language           string         `json:"language,omitempty"`
		Priority           string         `json:"priority,omitempty"`
		Location           string         `json:"location,omitempty"`
		Latitude           float64        `json:"latitude,omitempty"`
		Longitude          float64        `json:"longitude,omitempty"`
		DistanceFromTarget float64        `json:"-"`
		Visibility         uint8          `json:"visibility,omitempty"`
		Object             *Object        `json:"object"`
		Target             *Object        `json:"target,omitempty"`
		Instrument         *Object        `json:"instrument,omitempty"`
		Participant        []*Participant `json:"participant,omitempty"`
		Common
	}

	// EventWithIDs holds the event structure with the added account and application ids
	EventWithIDs struct {
		AccAppIDs
		Event
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

	// EventsResponse represents the common structure for responses which contains events
	EventsResponse struct {
		Events      []*Event                    `json:"events"`
		Users       map[string]*ApplicationUser `json:"users"`
		EventsCount int                         `json:"events_count"`
		UsersCount  int                         `json:"users_count"`
	}

	// EventsResponseWithUnread represents the common structure for responses which contains events and have an unread count
	EventsResponseWithUnread struct {
		EventsResponse
		UnreadCount int `json:"unread_events_count"`
	}

	// AutoConnectSocialFriends holds the informatin that we have for auto-connecting social frieds
	AutoConnectSocialFriends struct {
		User              *ApplicationUserWithIDs `json:"user"`
		Type              string                  `json:"type"`
		OurStoredUsersIDs []*ApplicationUser      `json:"our_stored_users_ids"`
	}

	// SocialConnection holds the social connection information
	SocialConnection struct {
		User             *ApplicationUserWithIDs `json:"user"`
		Platform         string                  `json:"platform"`
		Type             string                  `json:"type"`
		SocialFriendsIDs []string                `json:"social_friends_ids"`
	}

	// ErrorResponse holds the structure of an error what's reported back to the user
	ErrorResponse struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url,omitempty"`
	}

	// ErrorsResponse holds the structure for multiple errors that are reported back to the user
	ErrorsResponse struct {
		Errors []ErrorResponse `json:"errors"`
	}
)

const (
	// EventPrivate flags that the event is private
	EventPrivate = 10

	// EventConnections flags that the event is shared with the connetions of the user
	EventConnections = 20

	// EventPublic flags that the event is public
	EventPublic = 30

	// EventGlobal flags that the event is public and visibile in the WHOLE app (use it with consideration)
	EventGlobal = 40
)

var (
	// Because I don't know another way to have optional json values and still have values

	// PTrue is a pointers to true
	PTrue *bool

	// PFalse is a pointers to false
	PFalse *bool
)

func init() {
	tr := true
	fl := false
	PTrue = &tr
	PFalse = &fl
}

func (e SortableEventsByDistance) Len() int      { return len(e) }
func (e SortableEventsByDistance) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e SortableEventsByDistance) Less(i, j int) bool {
	return e[i].DistanceFromTarget < e[j].DistanceFromTarget
}
