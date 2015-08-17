package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/entity"
)

func (p *pg) organizationUpdate(msg string) []errors.Error {
	updatedOrganization := entity.Organization{}
	err := json.Unmarshal([]byte(msg), &updatedOrganization)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingOrganization, er := p.organization.Read(updatedOrganization.ID)
	if er != nil {
		return er
	}

	_, er = p.organization.Update(*existingOrganization, updatedOrganization, false)
	return er
}

func (p *pg) organizationDelete(msg string) []errors.Error {
	organization := &entity.Organization{}
	err := json.Unmarshal([]byte(msg), organization)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.organization.Delete(organization)
}
