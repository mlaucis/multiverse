package response

import "github.com/tapglue/backend/v02/entity"

// SanitizeAccountUser will sanitize the account user for usage via the API
func SanitizeAccountUser(user *entity.AccountUser) {
	user.Password = ""
}

// SanitizeApplicationUsers sanitize a slice of application users
func SanitizeApplicationUsers(users []*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}

// SanitizeApplicationUsersMap sanitizes a map of application users
func SanitizeApplicationUsersMap(users map[string]*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}
