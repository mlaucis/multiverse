// Package postgres implements the writer logic for writing data to postgres
package postgres

import (
	"net/http"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	postgresCore "github.com/tapglue/multiverse/v03/core/postgres"
	redisCore "github.com/tapglue/multiverse/v03/core/redis"
	ksis "github.com/tapglue/multiverse/v03/storage/kinesis"
	"github.com/tapglue/multiverse/v03/storage/postgres"
	"github.com/tapglue/multiverse/v03/writer"

	"github.com/garyburd/redigo/redis"
)

type pg struct {
	ksis            ksis.Client
	pg              postgres.Client
	red             core.Application
	organization    core.Organization
	member          core.Member
	application     core.Application
	applicationUser core.ApplicationUser
	connection      core.Connection
	event           core.Event
}

var errUnknownMessage = errors.New(http.StatusInternalServerError, 1, "unknown message retrieved", "", false)

func (p *pg) ProcessMessages(channelName, msg string) errors.Error {
	var ers []errors.Error
	switch channelName {
	case ksis.StreamAccountUpdate:
		ers = p.organizationUpdate(msg)
	case ksis.StreamAccountDelete:
		ers = p.organizationDelete(msg)
	case ksis.StreamAccountUserCreate:
		ers = p.memberCreate(msg)
	case ksis.StreamAccountUserUpdate:
		ers = p.memberUpdate(msg)
	case ksis.StreamAccountUserDelete:
		ers = p.memberDelete(msg)
	case ksis.StreamApplicationCreate:
		ers = p.applicationCreate(msg)
	case ksis.StreamApplicationUpdate:
		ers = p.applicationUpdate(msg)
	case ksis.StreamApplicationDelete:
		ers = p.applicationDelete(msg)
	case ksis.StreamApplicationUserUpdate:
		ers = p.applicationUserUpdate(msg)
	case ksis.StreamApplicationUserDelete:
		ers = p.applicationUserDelete(msg)
	case ksis.StreamConnectionCreate:
		ers = p.connectionCreate(msg)
	case ksis.StreamConnectionConfirm:
		ers = p.connectionConfirm(msg)
	case ksis.StreamConnectionUpdate:
		ers = p.connectionUpdate(msg)
	case ksis.StreamConnectionAutoConnect:
		ers = p.connectionAutoConnect(msg)
	case ksis.StreamConnectionSocialConnect:
		ers = p.connectionSocialConnect(msg)
	case ksis.StreamConnectionDelete:
		ers = p.connectionDelete(msg)
	case ksis.StreamEventCreate:
		ers = p.eventCreate(msg)
	case ksis.StreamEventUpdate:
		ers = p.eventUpdate(msg)
	case ksis.StreamEventDelete:
		ers = p.eventDelete(msg)
	default:
		return errUnknownMessage.UpdateInternalMessage(msg)
	}

	// TODO we should really do something with the error here, not just ignore it like this, maybe?
	if ers != nil {
		for idx := range ers {
			return ers[idx].UpdateInternalMessage(ers[idx].InternalErrorWithLocation() + "\t" + msg)
		}
	}

	return nil
}

// New will return a new PosgreSQL writer
func New(kinesis ksis.Client, pgsql postgres.Client, red *redis.Pool) writer.Writer {
	redisClient := redisCore.NewApplication(red)

	return &pg{
		ksis:            kinesis,
		pg:              pgsql,
		red:             redisClient,
		organization:    postgresCore.NewOrganization(pgsql),
		member:          postgresCore.NewMember(pgsql),
		application:     postgresCore.NewApplication(pgsql, redisClient),
		applicationUser: postgresCore.NewApplicationUser(pgsql),
		connection:      postgresCore.NewConnection(pgsql),
		event:           postgresCore.NewEvent(pgsql),
	}
}
