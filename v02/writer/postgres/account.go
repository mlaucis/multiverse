package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

func (p *pg) accountUpdate(msg string) []errors.Error {
	updatedAccount := entity.Account{}
	err := json.Unmarshal([]byte(msg), &updatedAccount)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingAccount, er := p.account.Read(updatedAccount.ID)
	if er != nil {
		return er
	}

	_, er = p.account.Update(*existingAccount, updatedAccount, false)
	return er
}

func (p *pg) accountDelete(msg string) []errors.Error {
	account := &entity.Account{}
	err := json.Unmarshal([]byte(msg), account)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.account.Delete(account)
}
