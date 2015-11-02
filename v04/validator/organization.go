package validator

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
)

const (
	orgNameMin = 3
	orgNameMax = 40

	orgDescriptionMin = 0
	orgDescriptionMax = 100
)

// CreateOrganization validates an account on create
func CreateOrganization(organization *entity.Organization) (errs []errors.Error) {
	if !StringLengthBetween(organization.Name, orgNameMin, orgNameMax) {
		errs = append(errs, errmsg.ErrAccountNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(organization.Description, orgDescriptionMin, orgDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize.SetCurrentLocation())
	}

	if organization.ID != 0 {
		errs = append(errs, errmsg.ErrOrgIDIsAlreadySet.SetCurrentLocation())
	}

	if organization.AuthToken != "" {
		errs = append(errs, errmsg.ErrOrgTokenAlreadySet.SetCurrentLocation())
	}

	return
}

// UpdateOrganization validates an account on update
func UpdateOrganization(existingOrg, updatedOrg *entity.Organization) (errs []errors.Error) {
	if updatedOrg.ID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero.SetCurrentLocation())
	}

	if !StringLengthBetween(updatedOrg.Name, orgNameMin, orgNameMax) {
		errs = append(errs, errmsg.ErrAccountNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(updatedOrg.Description, orgDescriptionMin, orgDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize.SetCurrentLocation())
	}

	return
}
