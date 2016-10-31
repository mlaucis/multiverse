package controller

import (
	"fmt"

	"github.com/tapglue/multiverse/service/member"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// MemberCreateFunc creates a member for the current Org.
type MemberCreateFunc func(
	org *v04_entity.Organization,
	m *member.Member,
) (*member.Member, error)

// MemberCreate creates a member for the current Org.
func MemberCreate(members member.Service) MemberCreateFunc {
	return func(
		currentOrg *v04_entity.Organization,
		m *member.Member,
	) (*member.Member, error) {
		return nil, fmt.Errorf("MemberCreate not implemented")
	}
}
