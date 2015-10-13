package postgres

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

func (p *pg) applicationUserUpdate(msg string) []errors.Error {
	updatedApplicationUser := entity.ApplicationUserWithIDs{}
	err := json.Unmarshal([]byte(msg), &updatedApplicationUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingApplicationUser, er := p.applicationUser.Read(updatedApplicationUser.OrgID, updatedApplicationUser.AppID, updatedApplicationUser.ID, false)
	if er != nil {
		return er
	}

	_, er = p.applicationUser.Update(updatedApplicationUser.OrgID, updatedApplicationUser.AppID, *existingApplicationUser, updatedApplicationUser.ApplicationUser, false)
	return er
}

func (p *pg) applicationUserDelete(msg string) []errors.Error {
	applicationUser := &entity.ApplicationUserWithIDs{}
	err := json.Unmarshal([]byte(msg), applicationUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.applicationUser.Delete(applicationUser.OrgID, applicationUser.AppID, applicationUser.ID)
}
