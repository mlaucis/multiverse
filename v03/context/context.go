package context

import (
	ctx "github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/v03/entity"
)

type (
	tokenType uint8

	// Context holds the context for the current request
	Context struct {
		ctx.Context

		TokenType         tokenType
		OrganizationID    int64
		MemberID          int64
		ApplicationID     int64
		ApplicationUserID uint64
		Organization      *entity.Organization
		Member            *entity.Member
		Application       *entity.Application
		ApplicationUser   *entity.ApplicationUser
	}
)

// Here we define the supported token types
const (
	TokenTypeUnknown tokenType = iota
	TokenTypeOrganization
	TokenTypeOrganizationID
	TokenTypeApplication
	TokenTypeBackend
)
