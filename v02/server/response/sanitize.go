/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package response

import "github.com/tapglue/backend/v02/entity"

func SanitizeApplicationUsers(users []*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}

func SanitizeApplicationUsersMap(users map[string]*entity.ApplicationUser) {
	for idx := range users {
		users[idx].Password = ""
		users[idx].Enabled = false
		users[idx].SocialIDs = map[string]string{}
		users[idx].Activated = false
		users[idx].Email = ""
		users[idx].CreatedAt, users[idx].UpdatedAt, users[idx].LastLogin, users[idx].LastRead = nil, nil, nil, nil
	}
}
