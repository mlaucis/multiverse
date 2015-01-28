package server

import (
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
)

func AddCorrectAccount() *entity.Account {
	savedAccount, err := core.AddAccount(correctAccount, true)
	if err != nil {
		panic(err)
	}

	return savedAccount
}
