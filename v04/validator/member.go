package validator

import (
	"strings"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/utils"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	storageHelper "github.com/tapglue/multiverse/v04/storage/helper"
)

const (
	memberNameMin     = 2
	memberNameMax     = 40
	memberPasswordMin = 4
	memberPasswordMax = 60
)

// CreateMember validates an account user on create
func CreateMember(datastore core.Member, member *entity.Member) (errs []errors.Error) {
	if !StringLengthBetween(member.FirstName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberFirstNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(member.LastName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberLastNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(member.Username, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberUsernameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(member.Password, memberPasswordMin, memberPasswordMax) {
		errs = append(errs, errmsg.ErrMemberPasswordSize.SetCurrentLocation())
	}

	if member.OrgID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero.SetCurrentLocation())
	}

	if member.Email == "" || !IsValidEmail(member.Email) {
		errs = append(errs, errmsg.ErrMemberEmailInvalid.SetCurrentLocation())
	}

	if isDuplicate, err := DuplicateMemberEmail(datastore, member.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation())
		} else {
			errs = append(errs, err...)
		}
	}

	if isDuplicate, err := DuplicateMemberUsername(datastore, member.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation())
		} else {
			errs = append(errs, err...)
		}
	}

	return
}

// UpdateMember validates an account user on update
func UpdateMember(datastore core.Member, existingMember, updatedMember *entity.Member) (errs []errors.Error) {
	if !StringLengthBetween(updatedMember.FirstName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberFirstNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(updatedMember.LastName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberLastNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(updatedMember.Username, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberUsernameSize.SetCurrentLocation())
	}

	if updatedMember.Password != "" {
		if !StringLengthBetween(updatedMember.Password, memberPasswordMin, memberPasswordMax) {
			errs = append(errs, errmsg.ErrMemberPasswordSize.SetCurrentLocation())
		}
	}

	if updatedMember.Email == "" || !IsValidEmail(updatedMember.Email) {
		errs = append(errs, errmsg.ErrMemberEmailInvalid.SetCurrentLocation())
	}

	if existingMember.Email != updatedMember.Email {
		if isDuplicate, err := DuplicateMemberEmail(datastore, updatedMember.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation())
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if existingMember.Username != updatedMember.Username {
		if isDuplicate, err := DuplicateMemberUsername(datastore, updatedMember.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation())
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return
}

// MemberCredentialsValid checks is a certain member has the right credentials
func MemberCredentialsValid(password string, member *entity.Member) (errs []errors.Error) {
	pass, err := utils.Base64Decode(member.Password)
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage("invalid password parts").SetCurrentLocation()}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.ErrAuthPasswordMismatch.SetCurrentLocation()}
	}

	return
}

// DuplicateMemberEmail checks if the user e-mail is duplicate within the provided account
func DuplicateMemberEmail(datastore core.Member, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation()}
		}
	}

	return false, nil
}

// DuplicateMemberUsername checks if the username is duplicate within the provided account
func DuplicateMemberUsername(datastore core.Member, username string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByUsername(username); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation()}
		}
	}

	return false, nil
}
