package response

import "github.com/tapglue/multiverse/v04/entity"

// SanitizeMember will sanitize the member for usage via the API
func SanitizeMember(member *entity.Member) {
	member.Password = ""
	member.Deleted = nil
}

// SanitizeApplicationUser will sanitize the application user for usage via the API
func SanitizeApplicationUser(user *entity.ApplicationUser) {
	user.Password = ""
	user.Deleted = nil
}

// SanitizeApplicationUsers sanitize a slice of application users
func SanitizeApplicationUsers(users []*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Deleted = nil
		users[idx].SessionToken = ""
		users[idx].FriendCount, users[idx].FollowerCount, users[idx].FollowedCount = nil, nil, nil
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}

// SanitizeApplicationUsersMap sanitizes a map of application users
func SanitizeApplicationUsersMap(users map[string]*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Deleted = nil
		users[idx].SessionToken = ""
		users[idx].FriendCount, users[idx].FollowerCount, users[idx].FollowedCount = nil, nil, nil
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}
