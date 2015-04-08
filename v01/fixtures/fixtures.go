/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 * @author Florin Patan <florinpatan@gmail.com>
 */

package fixtures

import "github.com/tapglue/backend/v01/entity"

// Definition of test data
var (
	EmtpyAccount     = entity.Account{}
	EmtpyAccountUser = entity.AccountUser{}
	EmptyApplication = entity.Application{}
	EmptyUser        = entity.User{}
	EmptyEvent       = entity.Event{}
	CorrectAccount   = entity.Account{
		Name:        "Demo",
		Description: "This is a demo account",
		Common: entity.Common{
			Enabled: true,
		},
	}
	CorrectAccountBig = entity.Account{
		Name: "Demozdr;aryprawurpaiw;ayeaasjhdsakjdlksajdlsakjlsakdjsalkdjlkasfja;sjflsakfaf[wor3pouarjlkfhzslkfhasfha;fha;kfhaslkgjas;lfjajsdhals;jfasljfhals;fja;skfhas;lfjas;kfhaslkghas;kghaslkghaslfhdlsakjdaslfasjas;lgjsaljgajgasjgasgas;k'saldksa;gosaeugauypaotyaptua;otyqpotyapyrqyrapytalktypawrpauadadasasdads",
		Common: entity.Common{
			Enabled: true,
		},
	}
	CorrectAccountUser = entity.AccountUser{
		UserCommon: entity.UserCommon{
			Username:  "Demo User",
			FirstName: "First name",
			LastName:  "Last Name",
			Password:  "iamsecure..not",
			Email:     "d@m.o",
		},
	}
	CorrectApplication = entity.Application{
		Name:        "Demo App",
		Description: "This is the best application",
		URL:         "http://tapglue.com",
	}
	CorrectUser = entity.User{
		UserCommon: entity.UserCommon{
			Username:  "demouser",
			FirstName: "Demo",
			LastName:  "User",
			Password:  "password",
			Email:     "user@tapglue.com",
			URL:       "http://tapglue.com/users/1/demouser",
		},
		Common: entity.Common{
			Image: []*entity.Image{
				{
					URL: "http://images.tapglue.com/1/demouser/profile.jpg",
				},
			},
			Metadata: "{\"customData\":\"customValue\"}",
		},
	}
	CorrectConnection = entity.Connection{
		Common: entity.Common{
			Enabled: true,
		},
	}
	CorrectEvent = entity.Event{
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
