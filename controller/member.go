package controller

import (
	"fmt"

	"github.com/tapglue/multiverse/platform/generate"

	"github.com/tapglue/multiverse/service/member"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// MemberCreateFunc creates a member for the current Org.
type MemberCreateFunc func(
	currentOrg *v04_entity.Organization,
	m *member.Member,
) (*member.Member, error)

// MemberCreate creates a member for the current Org.
func MemberCreate(members member.Service) MemberCreateFunc {
	return func(
		currentOrg *v04_entity.Organization,
		m *member.Member,
	) (*member.Member, error) {
		if err := m.Validate(); err != nil {
			return nil, wrapError(ErrInvalidEntity, "%s", err)
		}

		if err := constrainMemberEmailUnique(members, m.Email); err != nil {
			return nil, err
		}

		if err := constrainMemberUsernameUnique(members, m.Username); err != nil {
			return nil, err
		}

		pw, err := passwordSecure(m.Password)
		if err != nil {
			return nil, err
		}

		m.Enabled = true
		m.OrgID = uint64(currentOrg.ID)
		m.Password = pw
		m.PublicOrgID = currentOrg.PublicID

		m.PublicID, err = generate.UUID()
		if err != nil {
			return nil, err
		}

		m, err = members.Put(member.NamespaceDefault, m)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("MemberCreate not implemented")
	}
}

func constrainMemberEmailUnique(members member.Service, email string) error {
	ms, err := members.Query(member.NamespaceDefault, member.QueryOpts{
		Enabled: &defaultEnabled,
		Emails: []string{
			email,
		},
	})
	if err != nil {
		return err
	}

	if len(ms) > 0 {
		return wrapError(ErrInvalidEntity, "email already in use")
	}

	return nil
}

func constrainMemberUsernameUnique(members member.Service, username string) error {
	ms, err := members.Query(member.NamespaceDefault, member.QueryOpts{
		Enabled: &defaultEnabled,
		Usernames: []string{
			username,
		},
	})
	if err != nil {
		return err
	}

	if len(ms) > 0 {
		return wrapError(ErrInvalidEntity, "username already in use")
	}

	return nil
}
