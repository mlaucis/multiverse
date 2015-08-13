package response

import "github.com/tapglue/backend/v03/entity"

// SanitizeAccountUser will sanitize the account user for usage via the API
func SanitizeAccountUser(user *entity.AccountUser) {
	user.Password = ""
	user.Deleted = nil
}

// SanitizeApplicationUsers sanitize a slice of application users
func SanitizeApplicationUsers(users []*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].Activated = false
		users[idx].Deleted = nil
		users[idx].Email = ""
		users[idx].SessionToken = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}

// SanitizeApplicationUsersMap sanitizes a map of application users
func SanitizeApplicationUsersMap(users map[string]*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].Deleted = nil
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].SessionToken = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}
