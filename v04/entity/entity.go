//Package entity provides all the entities needed by the app to interact with the database
package entity

import (
	"fmt"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/errmsg"
)

const (
	rateLimitProduction = 20000
	rateLimitStaging    = 100
)

type (
	// OrgAppIDs holds the account and application IDs
	OrgAppIDs struct {
		OrgID         int64  `json:"acc_id,omitempty"`
		AppID         int64  `json:"app_id,omitempty"`
		PublicOrgID   string `json:"pub_acc_id,omitempty"`
		PublicAppID   string `json:"pub_app_id,omitempty"`
		CurrentUserID uint64 `json:"current_user_id,omitempty"`
	}

	// Common holds common used fields
	Common struct {
		Metadata  interface{}      `json:"metadata,omitempty"`
		Images    map[string]Image `json:"images,omitempty"`
		CreatedAt *time.Time       `json:"created_at,omitempty"`
		UpdatedAt *time.Time       `json:"updated_at,omitempty"`
		Enabled   bool             `json:"enabled"`
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
		ID           interface{}       `json:"id"`
		Type         string            `json:"type"`
		URL          string            `json:"url,omitempty"`
		DisplayNames map[string]string `json:"display_names,omitempty"` // ["en"=>"article", "de"=>"artikel"]
	}

	// Participant structure
	Participant struct {
		ID     interface{}      `json:"id"`
		URL    string           `json:"url,omitempty"`
		Images map[string]Image `json:"images,omitempty"`
	}

	// Organization structure
	Organization struct {
		ID           int64          `json:"-"`
		PublicID     string         `json:"id"`
		Name         string         `json:"name"`
		Description  string         `json:"description"`
		AuthToken    string         `json:"token"`
		Members      []*Member      `json:"-"`
		Applications []*Application `json:"-"`
		Common
	}

	// Member structure
	Member struct {
		ID              int64  `json:"-"`
		OrgID           int64  `json:"-"`
		PublicID        string `json:"id"`
		PublicAccountID string `json:"account_id"`
		SessionToken    string `json:"-"`
		UserCommon
		Common
	}

	// ApplicationUser structure
	ApplicationUser struct {
		ID                    uint64              `json:"id"`
		CustomID              string              `json:"custom_id,omitempty"`
		SocialIDs             map[string]string   `json:"social_ids,omitempty"`
		SocialConnectionsIDs  map[string][]string `json:"social_connections_ids,omitempty"`
		SocialConnectionType  ConnectionTypeType  `json:"connection_type,omitempty"`
		SocialConnectionState ConnectionStateType `json:"connection_state,omitempty"`
		DeviceIDs             []string            `json:"device_ids,omitempty"`
		Events                []*Event            `json:"events,omitempty"`
		Connections           []*ApplicationUser  `json:"connections,omitempty"`
		LastRead              *time.Time          `json:"-"`
		FriendCount           *int64              `json:"friend_count,omitempty"`
		FollowerCount         *int64              `json:"follower_count,omitempty"`
		FollowedCount         *int64              `json:"followed_count,omitempty"`
		Relation
		UserCommon
		Common
	}

	// ApplicationUserWithIDs holds the application user structure with the added account and application ids
	ApplicationUserWithIDs struct {
		OrgAppIDs
		ApplicationUser
	}

	// Connection structure holds the connections of the users
	Connection struct {
		UserFromID uint64              `json:"user_from_id"`
		UserToID   uint64              `json:"user_to_id"`
		Type       ConnectionTypeType  `json:"type"`
		State      ConnectionStateType `json:"state"`
		Enabled    *bool               `json:"enabled,omitempty"`
		CreatedAt  *time.Time          `json:"created_at,omitempty"`
		UpdatedAt  *time.Time          `json:"updated_at,omitempty"`
	}

	// ConnectionWithIDs holds the connection structure with the added account and application ids
	ConnectionWithIDs struct {
		OrgAppIDs
		Connection
	}

	// Relation holds the relation between two users
	Relation struct {
		IsFriend   *bool `json:"is_friend,omitempty"`
		IsFollower *bool `json:"is_follower,omitempty"`
		IsFollowed *bool `json:"is_followed,omitempty"`
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
		ID                 uint64        `json:"id"`
		UserID             uint64        `json:"user_id"`
		Type               string        `json:"type"`
		Language           string        `json:"language,omitempty"`
		Priority           string        `json:"priority,omitempty"`
		Location           string        `json:"location,omitempty"`
		Latitude           float64       `json:"latitude,omitempty"`
		Longitude          float64       `json:"longitude,omitempty"`
		DistanceFromTarget float64       `json:"-"`
		Visibility         uint8         `json:"visibility,omitempty"`
		Object             *Object       `json:"object"`
		Target             *Object       `json:"target,omitempty"`
		Instrument         *Object       `json:"instrument,omitempty"`
		Participant        []Participant `json:"participant,omitempty"`
		Common
	}

	// EventWithIDs holds the event structure with the added account and application ids
	EventWithIDs struct {
		OrgAppIDs
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

	// CreateSocialConnectionRequest is used by the client when requesting creation of a social connection from a user
	CreateSocialConnectionRequest struct {
		PlatformUserID  string              `json:"platform_user_id"`
		SocialPlatform  string              `json:"platform"`
		ConnectionsIDs  []string            `json:"connection_ids"`
		ConnectionType  ConnectionTypeType  `json:"type"`
		ConnectionState ConnectionStateType `json:"state"`
	}

	// ConnectionsByStateResponse is used as a response to the API query to return the connections in a certain state
	ConnectionsByStateResponse struct {
		IncomingConnections      []*Connection      `json:"incoming"`
		OutgoingConnections      []*Connection      `json:"outgoing"`
		Users                    []*ApplicationUser `json:"users"`
		IncomingConnectionsCount int                `json:"incoming_connections_count"`
		OutgoingConnectionsCount int                `json:"outgoing_connections_count"`
		UsersCount               int                `json:"users_count"`
	}

	// ConnectionStateType represents the type of a connection state
	ConnectionStateType string

	// ConnectionTypeType represents the type of a connection type
	ConnectionTypeType string
)

// Application structure
type Application struct {
	ID           int64              `json:"-"`
	OrgID        int64              `json:"-"`
	PublicID     string             `json:"id"`
	PublicOrgID  string             `json:"account_id"`
	AuthToken    string             `json:"token"`
	BackendToken string             `json:"backend_token"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	URL          string             `json:"url"`
	InProduction bool               `json:"in_production"`
	Users        []*ApplicationUser `json:"-"`
	Common
}

// Limit returns the desired rate limit for an Application varies by production
// state.
func (a *Application) Limit() int64 {
	if a.InProduction {
		return rateLimitProduction
	}

	return rateLimitStaging
}

// Namespace returrs the prefix for bucketing entities by Application.
func (a *Application) Namespace() string {
	return fmt.Sprintf("app_%d_%d", a.ID, a.OrgID)
}

const (
	// EventPrivate flags that the event is private
	EventPrivate = 10

	// EventConnections flags that the event is shared with the connetions of the user
	EventConnections = 20

	// EventPublic flags that the event is public
	EventPublic = 30

	// EventGlobal flags that the event is public and visibile in the WHOLE app (use it with consideration)
	EventGlobal = 40

	// ConnectionTypeFollow is a friend connection
	ConnectionTypeFriend ConnectionTypeType = "friend"

	// ConnectionTypeFollow is a follow connection
	ConnectionTypeFollow ConnectionTypeType = "follow"

	// ConnectionStatePending is used for pending connections
	ConnectionStatePending ConnectionStateType = "pending"

	// ConnectionStateConfirmed is used for accepted connections
	ConnectionStateConfirmed ConnectionStateType = "confirmed"

	// ConnectionStateRejected is used for rejected connections
	ConnectionStateRejected ConnectionStateType = "rejected"
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

// IsValidType will check if the current connection type is valid
func (c *Connection) IsValidType() bool {
	return c.Type == ConnectionTypeFollow ||
		c.Type == ConnectionTypeFriend
}

// IsValidState will check if the current connection state is valid
func (c *Connection) IsValidState() bool {
	return IsValidConectionState(c.State)
}

// TransferState will take care of transfering the connection state to the new state
func (c *Connection) TransferState(newState ConnectionStateType, issuerUserID uint64) []errors.Error {
	if !IsValidConectionState(newState) {
		return []errors.Error{errmsg.ErrConnectionStateInvalid.
			UpdateInternalMessage("got connection state: " + string(newState)).
			SetCurrentLocation()}
	}

	if c.State == "" {
		return c.transferEmpty(newState, issuerUserID)
	}

	if c.State == ConnectionStatePending {
		return c.transferPending(newState, issuerUserID)
	}

	if c.State == ConnectionStateConfirmed {
		return c.transferConfirmed(newState, issuerUserID)
	}

	if c.State == ConnectionStateRejected {
		return c.transferRejected(newState, issuerUserID)
	}

	return []errors.Error{errmsg.ErrConnectionStateInvalid.
		UpdateInternalMessage("failed to transfer connection to new state: " + string(newState) + " from state: " + string(c.State)).
		SetCurrentLocation()}
}

// IsValidConnectionState will check if the desired connection state is valid
func IsValidConectionState(state ConnectionStateType) bool {
	return state == ConnectionStatePending ||
		state == ConnectionStateConfirmed ||
		state == ConnectionStateRejected
}

func (c *Connection) transferEmpty(newState ConnectionStateType, issuerUserID uint64) []errors.Error {
	c.State = newState
	return nil
}

func (c *Connection) transferPending(newState ConnectionStateType, issuerUserID uint64) []errors.Error {
	if newState != ConnectionStateConfirmed &&
		newState != ConnectionStateRejected {
		return []errors.Error{errmsg.ErrConnectionStateNotAllowed.
			UpdateInternalMessage("failed to transfer connection to new state: " + string(newState) + " from state: " + string(c.State)).
			SetCurrentLocation()}
	}

	if issuerUserID != c.UserToID {
		return []errors.Error{errmsg.ErrConnectionStateTransferNotAllowed.SetCurrentLocation()}
	}

	c.State = newState
	return nil
}

func (c *Connection) transferConfirmed(newState ConnectionStateType, issuerUserID uint64) []errors.Error {
	if newState != ConnectionStateRejected {
		return []errors.Error{errmsg.ErrConnectionStateNotAllowed.
			UpdateInternalMessage("failed to transfer connection to new state: " + string(newState) + " from state: " + string(c.State)).
			SetCurrentLocation()}
	}

	if issuerUserID != c.UserToID {
		return []errors.Error{errmsg.ErrConnectionStateTransferNotAllowed.SetCurrentLocation()}
	}

	c.State = newState
	return nil
}

func (c *Connection) transferRejected(newState ConnectionStateType, issuerUserID uint64) []errors.Error {
	if newState != ConnectionStateConfirmed {
		return []errors.Error{errmsg.ErrConnectionStateNotAllowed.
			UpdateInternalMessage("failed to transfer connection to new state: " + string(newState) + " from state: " + string(c.State)).
			SetCurrentLocation()}
	}

	if issuerUserID != c.UserToID {
		return []errors.Error{errmsg.ErrConnectionStateTransferNotAllowed.SetCurrentLocation()}
	}

	c.State = newState
	return nil
}
