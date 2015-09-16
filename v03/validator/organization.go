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
		errs = append(errs, errmsg.ErrAccountNameSize)
	}

	if !StringLengthBetween(organization.Description, orgDescriptionMin, orgDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize)
	}

	if organization.ID != 0 {
		errs = append(errs, errmsg.ErrOrgIDIsAlreadySet)
	}

	if organization.AuthToken != "" {
		errs = append(errs, errmsg.ErrOrgTokenAlreadySet)
	}

	if len(organization.Images) > 0 {
		if !checkImages(organization.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	return
}

// UpdateOrganization validates an account on update
func UpdateOrganization(existingOrg, updatedOrg *entity.Organization) (errs []errors.Error) {
	if updatedOrg.ID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero)
	}

	if !StringLengthBetween(updatedOrg.Name, orgNameMin, orgNameMax) {
		errs = append(errs, errmsg.ErrAccountNameSize)
	}

	if !StringLengthBetween(updatedOrg.Description, orgDescriptionMin, orgDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize)
	}

	if len(updatedOrg.Images) > 0 {
		if !checkImages(updatedOrg.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	return
}
