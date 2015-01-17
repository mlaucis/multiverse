package server

import (
	"github.com/tapglue/backend/entity"
	"github.com/tapglue/backend/mysql"
)

func AddCorrectAccount() *entity.Account {

	var account = &entity.Account{
		Name: "Demo",
	}

	savedAccount, err := mysql.AddAccount(account)
	if err != nil {
		panic(err)
	}

	return savedAccount
}
