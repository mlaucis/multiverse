package event

import (
	"github.com/tapglue/multiverse/errors"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type nopService struct{}

// NewNopService returns a nop implementation of StrangleService.
func NewNopService() StrangleService {
	return &nopService{}
}

func (s *nopService) Create(
	orgID, appID int64,
	userID uint64,
	event *v04_entity.Event,
) []errors.Error {
	return nil
}

func (s *nopService) Delete(
	orgID, appID int64,
	userID, eventID uint64,
) []errors.Error {
	return nil
}

func (s *nopService) ListAll(
	orgID, appID int64,
	condition v04_core.EventCondition,
) ([]*v04_entity.Event, []errors.Error) {
	return nil, nil
}
