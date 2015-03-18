/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import "github.com/tapglue/backend/core/entity"

// Definition of test data
var (
	emtpyAccount     = &entity.Account{}
	emtpyAccountUser = &entity.AccountUser{}
	emptyApplication = &entity.Application{}
	emptyUser        = &entity.User{}
	emptyEvent       = &entity.Event{}
	correctAccount   = &entity.Account{
		Name:        "Demo",
		Description: "This is a demo account",
		Common: entity.Common{
			Enabled: true,
		},
	}
	correctAccountBig = &entity.Account{
		Name: "Demozdr;aryprawurpaiw;ayeaasjhdsakjdlksajdlsakjlsakdjsalkdjlkasfja;sjflsakfaf[wor3pouarjlkfhzslkfhasfha;fha;kfhaslkgjas;lfjajsdhals;jfasljfhals;fja;skfhas;lfjas;kfhaslkghas;kghaslkghaslfhdlsakjdaslfasjas;lgjsaljgajgasjgasgas;k'saldksa;gosaeugauypaotyaptua;otyqpotyapyrqyrapytalktypawrpauadadasasdads",
		Common: entity.Common{
			Enabled: true,
		},
	}
	correctAccountUser = &entity.AccountUser{
		UserCommon: entity.UserCommon{
			Username:  "Demo User",
			FirstName: "First name",
			LastName:  "Last Name",
			Password:  "iamsecure..not",
			Email:     "d@m.o",
		},
	}
	correctApplication = &entity.Application{
		Name:        "Demo App",
		Description: "This is the best application",
		URL:         "http://app.co",
	}
	correctUser = &entity.User{
		AuthToken: "userToken123",
		UserCommon: entity.UserCommon{
			Username:  "Demo User",
			FirstName: "Florin",
			LastName:  "Patan",
			Password:  "password",
			Email:     "d@m.o",
			URL:       "http://link.to/userToken123",
		},
		Common: entity.Common{
			Image: []*entity.Image{
				{
					URL: "http://link.to/userthumbnail.jpg",
				},
			},
			Metadata: "{\"more\":\"data\"}",
		},
	}
	correctConnection = &entity.Connection{
		Common: entity.Common{
			Enabled: true,
		},
	}
	correctEvent = &entity.Event{
		Verb:     "like",
		Language: "en",
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
