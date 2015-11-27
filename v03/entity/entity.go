//Package entity provides all the entities needed by the app to interact with the database
package entity

import (
	"encoding/json"
	"strconv"
	"time"
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
		ID                   uint64              `json:"id"`
		CustomID             string              `json:"custom_id,omitempty"`
		SocialIDs            map[string]string   `json:"social_ids,omitempty"`
		SocialConnectionsIDs map[string][]string `json:"social_connections_ids,omitempty"`
		SocialConnectionType string              `json:"connection_type,omitempty"`
		DeviceIDs            []string            `json:"device_ids,omitempty"`
		Events               []*Event            `json:"events,omitempty"`
		Connections          []*ApplicationUser  `json:"connections,omitempty"`
		LastRead             *time.Time          `json:"-"`
		FriendCount          *int64              `json:"friend_count,omitempty"`
		FollowerCount        *int64              `json:"follower_count,omitempty"`
		FollowedCount        *int64              `json:"followed_count,omitempty"`
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
		UserFromID  uint64     `json:"user_from_id"`
		UserToID    uint64     `json:"user_to_id"`
		Type        string     `json:"type"`
		State       string     `json:"state,omitempty"`
		ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
		Common
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
		Object             *Object       `json:"object"`
		Target             *Object       `json:"target,omitempty"`
		Instrument         *Object       `json:"instrument,omitempty"`
		Participant        []Participant `json:"participant,omitempty"`
		Common
	}

	// PresentationEvent holds the struct used to represent the event to the outside world
	PresentationEvent struct {
		IDString     string `json:"id_string"`
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

	// ConnectionTypeFollow is a friend connection
	ConnectionTypeFriend = "friend"

	// ConnectionTypeFollow is a follow connection
	ConnectionTypeFollow = "follow"
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

func (e *PresentationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IDString     string `json:"id_string"`
		UserIDString string `json:"user_id_string"`
		*Event
	}{
		IDString:     strconv.FormatUint(e.ID, 10),
		UserIDString: strconv.FormatUint(e.UserID, 10),
		Event:        e.Event,
	})
}

func (e SortableEventsByDistance) Len() int      { return len(e) }
func (e SortableEventsByDistance) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e SortableEventsByDistance) Less(i, j int) bool {
	return e[i].DistanceFromTarget < e[j].DistanceFromTarget
}
