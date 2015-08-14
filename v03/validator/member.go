package validator

import (
	"strings"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	storageHelper "github.com/tapglue/backend/v03/storage/helper"
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
		errs = append(errs, errmsg.ErrMemberFirstNameSize)
	}

	if !StringLengthBetween(member.LastName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberLastNameSize)
	}

	if !StringLengthBetween(member.Username, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberUsernameSize)
	}

	if !StringLengthBetween(member.Password, memberPasswordMin, memberPasswordMax) {
		errs = append(errs, errmsg.ErrMemberPasswordSize)
	}

	if member.OrgID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero)
	}

	if member.Email == "" || !IsValidEmail(member.Email) {
		errs = append(errs, errmsg.ErrMemberEmailInvalid)
	}

	if member.URL != "" && !IsValidURL(member.URL, false) {
		errs = append(errs, errmsg.ErrMemberURLInvalid)
	}

	if len(member.Images) > 0 {
		if !checkImages(member.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateMemberEmail(datastore, member.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists)
		} else {
			errs = append(errs, err...)
		}
	}

	if isDuplicate, err := DuplicateMemberUsername(datastore, member.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserUsernameInUse)
		} else {
			errs = append(errs, err...)
		}
	}

	return
}

// UpdateMember validates an account user on update
func UpdateMember(datastore core.Member, existingMember, updatedMember *entity.Member) (errs []errors.Error) {
	if !StringLengthBetween(updatedMember.FirstName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberFirstNameSize)
	}

	if !StringLengthBetween(updatedMember.LastName, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberLastNameSize)
	}

	if !StringLengthBetween(updatedMember.Username, memberNameMin, memberNameMax) {
		errs = append(errs, errmsg.ErrMemberUsernameSize)
	}

	if updatedMember.Password != "" {
		if !StringLengthBetween(updatedMember.Password, memberPasswordMin, memberPasswordMax) {
			errs = append(errs, errmsg.ErrMemberPasswordSize)
		}
	}

	if updatedMember.Email == "" || !IsValidEmail(updatedMember.Email) {
		errs = append(errs, errmsg.ErrMemberEmailInvalid)
	}

	if updatedMember.URL != "" && !IsValidURL(updatedMember.URL, true) {
		errs = append(errs, errmsg.ErrMemberURLInvalid)
	}

	if len(updatedMember.Images) > 0 {
		if !checkImages(updatedMember.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	if existingMember.Email != updatedMember.Email {
		if isDuplicate, err := DuplicateMemberEmail(datastore, updatedMember.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if existingMember.Username != updatedMember.Username {
		if isDuplicate, err := DuplicateMemberUsername(datastore, updatedMember.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserUsernameInUse)
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
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage("invalid password parts")}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.ErrAuthPasswordMismatch}
	}

	return
}

// DuplicateMemberEmail checks if the user e-mail is duplicate within the provided account
func DuplicateMemberEmail(datastore core.Member, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.ErrApplicationUserEmailAlreadyExists}
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
			return true, []errors.Error{errmsg.ErrApplicationUserUsernameInUse}
		}
	}

	return false, nil
}
