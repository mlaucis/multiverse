package user

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type nopService struct{}

// NewNopService returns a nop implmentation of StrangleService.
func NewNopService() StrangleService {
	return &nopService{}
}

func (s *nopService) FilterByEmail(
	orgID, appID int64,
	emails []string,
) ([]*v04_entity.ApplicationUser, []errors.Error) {
	return nil, nil
}

func (s *nopService) FilterBySocialIDs(
	orgID, appID int64,
	platform string,
	ids []string,
) ([]*v04_entity.ApplicationUser, []errors.Error) {
	return nil, nil
}

func (s *nopService) FindBySession(
	orgID, appID int64,
	key string,
) (*v04_entity.ApplicationUser, []errors.Error) {
	return nil, nil
}

func (s *nopService) Read(
	orgID, appID int64,
	id uint64,
	stats bool,
) (*v04_entity.ApplicationUser, []errors.Error) {
	return nil, nil
}

func (s *nopService) UpdateLastRead(
	orgID, appID int64,
	id uint64,
) []errors.Error {
	return nil
}
