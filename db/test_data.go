package db

import "github.com/tapglue/backend/entity"

// Definition of test data
var (
	emtpyAccount     = &entity.Account{}
	emtpyAccountUser = &entity.AccountUser{}
	emptyApplication = &entity.Application{}
	emptyUser        = &entity.User{}
	emptySessions    = &entity.Session{}
	emptyEvent       = &entity.Event{}
	correctAccount   = &entity.Account{
		Name: "Demo",
	}
	correctAccountUser = &entity.AccountUser{
		Name:     "Demo User",
		Password: "iamsecure..not",
		Email:    "d@m.o",
	}
	correctApplication = &entity.Application{
		Name: "Demo App",
		Key:  "imanappkey12345",
	}
	correctUser = &entity.User{
		Token:        "userToken123",
		Username:     "Demo User",
		Name:         "Florin",
		Password:     "password",
		Email:        "d@m.o",
		URL:          "app://link.to/userToken123",
		ThumbnailURL: "http://link.to/userthumbnail.jpg",
		Provider:     "Native",
		Custom:       "{\"more\":\"data\"}",
	}
	correctSession = &entity.Session{
		AppID:     1,
		UserToken: "userToken123",
		Custom:    "{\"more\":\"data\"}",
	}
)
