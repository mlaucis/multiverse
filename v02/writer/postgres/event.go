package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

func (p *pg) eventCreate(msg string) []errors.Error {
	event := &entity.EventWithIDs{}
	err := json.Unmarshal([]byte(msg), event)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.event.Create(event.AccountID, event.ApplicationID, event.CurrentUserID, &event.Event, false)
	return er
}

func (p *pg) eventUpdate(msg string) []errors.Error {
	updatedEvent := entity.EventWithIDs{}
	err := json.Unmarshal([]byte(msg), &updatedEvent)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingEvent, er := p.event.Read(
		updatedEvent.AccountID,
		updatedEvent.ApplicationID,
		updatedEvent.CurrentUserID,
		updatedEvent.UserID,
		updatedEvent.ID)
	if er != nil {
		return er
	}

	_, er = p.event.Update(
		updatedEvent.AccountID,
		updatedEvent.ApplicationID,
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

	return p.event.Delete(event.AccountID, event.ApplicationID, event.CurrentUserID, &event.Event)
}
