//Package entity provides all the entities needed by the app to interact with the database
package entity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/errmsg"
)

// TypeTargetUser is the canonical identifier for internal user representations.
const TypeTargetUser = "tg_user"

// TypeEvent is the canonical identifier for internal events.
const (
	TypeEventFriend = "tg_friend"
	TypeEventFollow = "tg_follow"
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
		FirstName        string     `json:"first_name"`
		LastName         string     `json:"last_name"`
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

	// Application structure
	Application struct {
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

	// PresentationApplicationUser holds the struct used to represent the user to the outside world
	PresentationApplicationUser struct {
		IDString string `json:"id_string"`
		*ApplicationUser
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

	// PresentationConnection holds the struct used to represent the connection to the outside world
	PresentationConnection struct {
		UserFromIDString string `json:"user_from_id_string"`
		UserToIDString   string `json:"user_to_id_string"`
		*Connection
	}

	// Relation holds the relation between two users
	Relation struct {
		IsFriend   *bool `json:"is_friend,omitempty"`
		IsFollower *bool `json:"is_follower,omitempty"`
		IsFollowed *bool `json:"is_followed,omitempty"`
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
		Object             *Object       `json:"object,omitempty"`
		ObjectID           uint64        `json:"object_id"`
		Owned              bool          `json:"owned"`
		Target             *Object       `json:"target,omitempty"`
		Instrument         *Object       `json:"instrument,omitempty"`
		Participant        []Participant `json:"participant,omitempty"`
		Common
	}

	// PresentationEvent holds the struct used to represent the event to the outside world
	PresentationEvent struct {
		IDString     string `json:"id_string"`
		TGObjectID   string `json:"tg_object_id"`
		UserIDString string `json:"user_id_string"`
		*Event
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
		Events      []*PresentationEvent                    `json:"events"`
		Users       map[string]*PresentationApplicationUser `json:"users"`
		EventsCount int                                     `json:"events_count"`
		UsersCount  int                                     `json:"users_count"`
	}

	// EventsResponseWithUnread represents the common structure for responses which contains events and have an unread count
	EventsResponseWithUnread struct {
		EventsResponse
		UnreadCount int `json:"unread_events_count"`
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
		IncomingConnections      []*PresentationConnection      `json:"incoming"`
		OutgoingConnections      []*PresentationConnection      `json:"outgoing"`
		Users                    []*PresentationApplicationUser `json:"users"`
		IncomingConnectionsCount int                            `json:"incoming_connections_count"`
		OutgoingConnectionsCount int                            `json:"outgoing_connections_count"`
		UsersCount               int                            `json:"users_count"`
	}

	// ConnectionStateType represents the type of a connection state
	ConnectionStateType string

	// ConnectionTypeType represents the type of a connection type
	ConnectionTypeType string
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

	// ConnectionTypeFriend is a friend connection
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
	return fmt.Sprintf("app_%d_%d", a.OrgID, a.ID)
}

func (u *PresentationApplicationUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDString string `json:"id_string"`
		*ApplicationUser
	}{
		IDString:        strconv.FormatUint(u.ID, 10),
		ApplicationUser: u.ApplicationUser,
	})
}

func (conn *PresentationConnection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UserFromIDString string `json:"user_from_id_string"`
		UserToIDString   string `json:"user_to_id_string"`
		*Connection
	}{
		UserFromIDString: strconv.FormatUint(conn.UserFromID, 10),
		UserToIDString:   strconv.FormatUint(conn.UserToID, 10),
		Connection:       conn.Connection,
	})
}

func (p *PresentationConnection) UnmarshalJSON(raw []byte) error {
	f := struct {
		State            ConnectionStateType `json:"state"`
		Type             ConnectionTypeType  `json:"type"`
		UserFromID       uint64              `json:"user_from_id"`
		UserFromIDString string              `json:"user_from_id_string"`
		UserToID         uint64              `json:"user_to_id"`
		UserToIDString   string              `json:"user_to_id_string"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.Connection = &Connection{
		State:      f.State,
		Type:       f.Type,
		UserFromID: f.UserFromID,
		UserToID:   f.UserToID,
	}
	p.UserFromIDString = f.UserFromIDString
	p.UserToIDString = f.UserToIDString

	if f.UserFromID == 0 && f.UserFromIDString != "" {
		id, err := strconv.ParseUint(p.UserFromIDString, 10, 64)
		if err != nil {
			return err
		}

		p.Connection.UserFromID = id
	}

	if f.UserToID == 0 && f.UserToIDString != "" {
		id, err := strconv.ParseUint(p.UserToIDString, 10, 64)
		if err != nil {
			return err
		}

		p.Connection.UserToID = id
	}

	return nil
}

func (e *PresentationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDString     string `json:"id_string"`
		TGObjectID   string `json:"tg_object_id"`
		UserIDString string `json:"user_id_string"`
		*Event
	}{
		IDString:     strconv.FormatUint(e.ID, 10),
		TGObjectID:   strconv.FormatUint(e.Event.ObjectID, 10),
		UserIDString: strconv.FormatUint(e.UserID, 10),
		Event:        e.Event,
	})
}

func (e *PresentationEvent) UnmarshalJSON(raw []byte) error {
	r := struct {
		ID           uint64        `json:"id"`
		IDString     string        `json:"id_string"`
		UserID       uint64        `json:"user_id"`
		UserIDString string        `json:"user_id_string"`
		Type         string        `json:"type"`
		Language     string        `json:"language,omitempty"`
		Priority     string        `json:"priority,omitempty"`
		Location     string        `json:"location,omitempty"`
		Latitude     float64       `json:"latitude,omitempty"`
		Longitude    float64       `json:"longitude,omitempty"`
		Visibility   uint8         `json:"visibility,omitempty"`
		Object       *Object       `json:"object"`
		ObjectID     string        `json:"tg_object_id"`
		Target       *Object       `json:"target,omitempty"`
		Instrument   *Object       `json:"instrument,omitempty"`
		Participant  []Participant `json:"participant,omitempty"`
		Common
	}{}

	err := json.Unmarshal(raw, &r)
	if err != nil {
		return err
	}

	e.IDString = r.IDString
	e.TGObjectID = r.ObjectID
	e.UserIDString = r.UserIDString
	e.Event = &Event{
		ID:          r.ID,
		UserID:      r.UserID,
		Type:        r.Type,
		Language:    r.Language,
		Priority:    r.Priority,
		Location:    r.Location,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Visibility:  r.Visibility,
		Object:      r.Object,
		Target:      r.Target,
		Instrument:  r.Instrument,
		Participant: r.Participant,
		Common:      r.Common,
	}

	if r.ObjectID != "" {
		id, err := strconv.ParseUint(r.ObjectID, 10, 64)
		if err != nil {
			return err
		}

		e.Event.ObjectID = id
	}

	return nil
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
