package server

import (
	"github.com/tapglue/backend/db"
	"github.com/tapglue/backend/entity"
)

func AddCorrectAccount() *entity.Account {

	var account = &entity.Account{
		Name: "Demo",
	}

	savedAccount, err := db.AddAccount(account)
	if err != nil {
		panic(err)
	}

	return savedAccount
}
