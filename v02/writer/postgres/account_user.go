package postgres

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v02/entity"
)

func (p *pg) accountUserCreate(msg string) []errors.Error {
	accountUser := &entity.AccountUser{}
	err := json.Unmarshal([]byte(msg), accountUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.accountUser.Create(accountUser, false)
	return er
}

func (p *pg) accountUserUpdate(msg string) []errors.Error {
	updatedAccountUser := entity.AccountUser{}
	err := json.Unmarshal([]byte(msg), &updatedAccountUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingAccountUser, er := p.accountUser.Read(updatedAccountUser.AccountID, updatedAccountUser.ID)
	if er != nil {
		return er
	}

	_, er = p.accountUser.Update(*existingAccountUser, updatedAccountUser, false)
	return er
}

func (p *pg) accountUserDelete(msg string) []errors.Error {
	accountUser := &entity.AccountUser{}
	err := json.Unmarshal([]byte(msg), accountUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.accountUser.Delete(accountUser)
}
