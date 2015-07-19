package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

func (p *pg) applicationCreate(msg string) []errors.Error {
	application := &entity.Application{}
	err := json.Unmarshal([]byte(msg), application)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.application.Create(application, false)
	return er
}

func (p *pg) applicationUpdate(msg string) []errors.Error {
	updatedApplication := entity.Application{}
	err := json.Unmarshal([]byte(msg), &updatedApplication)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingApplication, er := p.application.Read(updatedApplication.AccountID, updatedApplication.ID)
	if er != nil {
		return nil
	}

	_, er = p.application.Update(*existingApplication, updatedApplication, false)
	return er
}

func (p *pg) applicationDelete(msg string) []errors.Error {
	application := &entity.Application{}
	err := json.Unmarshal([]byte(msg), application)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.application.Delete(application)
}
