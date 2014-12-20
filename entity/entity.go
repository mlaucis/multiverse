/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package entity

type (
	Application struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	User struct {
		AppID        uint64  `json:"app_id,omitempty"`
		Token        string  `json:"token,omitempty"`
		DisplayName  string  `json:"display_name,omitempty"`
		URL          string  `json:"url,omitempty"`
		ThumbnailUrl string  `json:"thumbnail_url,omitempty"`
		Custom       string  `json:"custom,omitempty"`
		CreatedAt    string  `json:"created_at,omitempty"`
		UpdatedAt    string  `json:"updated_at,omitempty"`
		LastLogin    string  `json:"last_login,omitempty"`
		Friends      []*User `json:"friends,omitempty"`
	}

	Event struct {
		ID           uint64 `json:"eventId"`
		AppID        uint64 `json:"app_id,omitempty"`
		UserID       uint64 `json:"user_id,omitempty"`
		EventType    string `json:"event_type"`
		ThumbnailUrl string `json:"thumbnail_url,omitempty"`
		ItemID       string `json:"item_id"`
		ItemName     string `json:"item_name,omitempty"`
		ItemURL      string `json:"item_url,omitempty"`
		CreatedAt    string `json:"created_at"`
		Custom       string `json:"custom,omitempty"`
		User         *User  `json:"user,omitempty"`
	}

	NewEvent struct {
		Verb string `json:"verb"`
		Event
		User
	}
)
