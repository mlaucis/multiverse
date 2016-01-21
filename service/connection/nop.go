package connection

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type nopService struct{}

// NewNopService returns a nop implementation of StrangleService.
func NewNopService() StrangleService {
	return &nopService{}
}

func (s *nopService) ConnectionsByState(
	orgID, appID int64,
	id uint64,
	state v04_entity.ConnectionStateType,
) ([]*v04_entity.Connection, []errors.Error) {
	return nil, nil
}

func (s *nopService) FriendsAndFollowingIDs(
	orgID, appID int64,
	id uint64,
) ([]uint64, []errors.Error) {
	return nil, nil
}

func (s *nopService) Relation(
	orgID, appID int64,
	from, to uint64,
) (*v04_entity.Relation, []errors.Error) {
	return nil, nil
}
