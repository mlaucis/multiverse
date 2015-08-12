package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/entity"
)

func (p *pg) applicationUserUpdate(msg string) []errors.Error {
	updatedApplicationUser := entity.ApplicationUserWithIDs{}
	err := json.Unmarshal([]byte(msg), &updatedApplicationUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingApplicationUser, er := p.applicationUser.Read(updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, updatedApplicationUser.ID)

	_, er = p.applicationUser.Update(updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, *existingApplicationUser, updatedApplicationUser.ApplicationUser, false)
	return er
}

func (p *pg) applicationUserDelete(msg string) []errors.Error {
	applicationUser := &entity.ApplicationUserWithIDs{}
	err := json.Unmarshal([]byte(msg), applicationUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.applicationUser.Delete(applicationUser.AccountID, applicationUser.ApplicationID, &applicationUser.ApplicationUser)
}
