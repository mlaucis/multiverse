package server

import (
	"github.com/tapglue/backend/aerospike"
	"github.com/tapglue/backend/entity"
)

func AddCorrectAccount() *entity.Account {
	savedAccount, err := aerospike.AddAccount(correctAccount, true)
	if err != nil {
		panic(err)
	}

	return savedAccount
}
