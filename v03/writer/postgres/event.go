package postgres

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

func (p *pg) eventCreate(msg string) []errors.Error {
	event := &entity.EventWithIDs{}
	err := json.Unmarshal([]byte(msg), event)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.event.Create(event.OrgID, event.AppID, event.CurrentUserID, &event.Event)
}

func (p *pg) eventUpdate(msg string) []errors.Error {
	updatedEvent := entity.EventWithIDs{}
	err := json.Unmarshal([]byte(msg), &updatedEvent)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingEvent, er := p.event.Read(
		updatedEvent.OrgID,
		updatedEvent.AppID,
		updatedEvent.UserID,
		updatedEvent.ID)
	if er != nil {
		return er
	}

	_, er = p.event.Update(
		updatedEvent.OrgID,
		updatedEvent.AppID,
		updatedEvent.CurrentUserID,
		*existingEvent,
		updatedEvent.Event,
		false)
	return er
}

func (p *pg) eventDelete(msg string) []errors.Error {
	event := &entity.EventWithIDs{}
	err := json.Unmarshal([]byte(msg), event)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.event.Delete(event.OrgID, event.AppID, event.CurrentUserID, event.ID)
}
