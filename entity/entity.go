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
		ID          uint64 `json:"id,omitempty"`
		AppID       uint64 `json:"appId,omitempty"`
		Token       string `json:"token"`
		DisplayName string `json:"displayName,omitempty"`
		URL         string `json:"url,omitempty"`
	}

	Event struct {
		ID        uint64 `json:"id"`
		AppID     uint64 `json:"appId,omitempty"`
		UserID    uint64 `json:"userId,omitempty"`
		EventType string `json:"eventType"`
		ItemID    string `json:"itemId"`
		ItemURL   string `json:"itemUrl,omitempty"`
		CreatedAt uint64 `json:"createdAt"`
	}

	NewEvent struct {
		Verb string `json:"verb"`
		Event
		User
	}
)
