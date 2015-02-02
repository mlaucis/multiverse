package mysql

import "github.com/tapglue/backend/core/entity"

// Definition of test data
var (
	emtpyAccount     = &entity.Account{}
	emtpyAccountUser = &entity.AccountUser{}
	emptyApplication = &entity.Application{}
	emptyUser        = &entity.User{}
	emptyEvent       = &entity.Event{}
	correctAccount   = &entity.Account{
		Name:    "Demo",
		Enabled: true,
	}
	correctAccountBig = &entity.Account{
		Name:    "Demozdr;aryprawurpaiw;ayeaasjhdsakjdlksajdlsakjlsakdjsalkdjlkasfja;sjflsakfaf[wor3pouarjlkfhzslkfhasfha;fha;kfhaslkgjas;lfjajsdhals;jfasljfhals;fja;skfhas;lfjas;kfhaslkghas;kghaslkghaslfhdlsakjdaslfasjas;lgjsaljgajgasjgasgas;k'saldksa;gosaeugauypaotyaptua;otyqpotyapyrqyrapytalktypawrpauadadasasdads",
		Enabled: true,
	}
	correctAccountUser = &entity.AccountUser{
		UserCommon: entity.UserCommon{
			DisplayName: "Demo User",
			Password:    "iamsecure..not",
			Email:       "d@m.o",
		},
	}
	correctApplication = &entity.Application{
		Name:      "Demo App",
		AuthToken: "imanappkey12345",
	}
	correctUser = &entity.User{
		AuthToken: "userToken123",
		UserCommon: entity.UserCommon{
			Username:  "Demo User",
			FirstName: "Florin",
			Password:  "password",
			Email:     "d@m.o",
			URL:       "app://link.to/userToken123",
		},
		Common: entity.Common{
			Image: []*entity.Image{
				&entity.Image{
					URL: "http://link.to/userthumbnail.jpg",
				},
			},
			Metadata: "{\"more\":\"data\"}",
		},
	}
	correctEvent = &entity.Event{
		ApplicationID: 1,
		UserID:        correctUser.ID,
		Object: &entity.Object{
			DisplayName: map[string]string{
				"en": "Event performed",
			},
		},
		Common: entity.Common{
			Metadata: "{\"more\":\"data\"}",
		},
	}
)