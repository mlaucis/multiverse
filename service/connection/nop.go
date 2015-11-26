package connection

import "github.com/tapglue/multiverse/errors"

type nopService struct{}

// NewNopService returns a nop implementation of StrangleService.
func NewNopService() StrangleService {
	return &nopService{}
}

func (s *nopService) FriendsAndFollowingIDs(
	orgID, appID int64,
	id uint64,
) ([]uint64, []errors.Error) {
	return nil, nil
}
