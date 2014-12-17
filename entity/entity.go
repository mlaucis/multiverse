/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package entity

type (
	Event struct {
		ID        string `json: "id"`
		EventType string `json: "eventType"`
		ItemID    string `json: "itemId"`
		ItemURL   string `json: "itemUrl,omitempty"`
		CreatedAt string `json: "createdAt"`
	}

	User struct {
		ID          string `json: "id,omitempty"`
		Token       string `json: "token"`
		DisplayName string `json: "displayName,omitempty"`
		URL         string `json: "url,omitempty"`
	}

	NewEvent struct {
		Verb string `json: "verb"`
		Event
		User
	}
)
