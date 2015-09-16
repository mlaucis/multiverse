package postgres

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

func (p *pg) memberCreate(msg string) []errors.Error {
	accountUser := &entity.Member{}
	err := json.Unmarshal([]byte(msg), accountUser)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.member.Create(accountUser, false)
	return er
}

func (p *pg) memberUpdate(msg string) []errors.Error {
	updatedMember := entity.Member{}
	err := json.Unmarshal([]byte(msg), &updatedMember)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingMember, er := p.member.Read(updatedMember.OrgID, updatedMember.ID)
	if er != nil {
		return er
	}

	_, er = p.member.Update(*existingMember, updatedMember, false)
	return er
}

func (p *pg) memberDelete(msg string) []errors.Error {
	member := &entity.Member{}
	err := json.Unmarshal([]byte(msg), member)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.member.Delete(member)
}
