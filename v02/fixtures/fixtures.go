/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 * @author Florin Patan <florinpatan@gmail.com>
 */

package fixtures

import . "github.com/tapglue/backend/v02/entity"

// Definition of test data
var (
	EmtpyAccount     = Account{}
	EmtpyAccountUser = AccountUser{}
	EmptyApplication = Application{}
	EmptyUser        = User{}
	EmptyEvent       = Event{}
	CorrectAccount   = Account{
		Name:        "Demo",
		Description: "This is a demo account",
		Common: Common{
			Enabled: true,
		},
	}
	CorrectAccountBig = Account{
		Name: "Demozdr;aryprawurpaiw;ayeaasjhdsakjdlksajdlsakjlsakdjsalkdjlkasfja;sjflsakfaf[wor3pouarjlkfhzslkfhasfha;fha;kfhaslkgjas;lfjajsdhals;jfasljfhals;fja;skfhas;lfjas;kfhaslkghas;kghaslkghaslfhdlsakjdaslfasjas;lgjsaljgajgasjgasgas;k'saldksa;gosaeugauypaotyaptua;otyqpotyapyrqyrapytalktypawrpauadadasasdads",
		Common: Common{
			Enabled: true,
		},
	}
	CorrectAccountUser = AccountUser{
		UserCommon: UserCommon{
			Username:  "Demo User",
			FirstName: "First name",
			LastName:  "Last Name",
			Password:  "iamsecure..not",
			Email:     "d@m.o",
		},
	}
	CorrectApplication = Application{
		Name:        "Demo App",
		Description: "This is the best application",
		URL:         "http://tapglue.com",
	}
	CorrectUser = User{
		UserCommon: UserCommon{
			Username:  "demouser",
			FirstName: "Demo",
			LastName:  "User",
			Password:  "password",
			Email:     "user@tapglue.com",
			URL:       "http://tapglue.com/users/1/demouser",
		},
		Common: Common{
			Image: []*Image{
				{
					URL: "http://images.tapglue.com/1/demouser/profile.jpg",
				},
			},
			Metadata: "{\"customData\":\"customValue\"}",
		},
	}
	CorrectConnection = Connection{
		Common: Common{
			Enabled: true,
		},
	}
	CorrectEvent = Event{
		Verb:     "like",
		Language: "en",
		Object: &Object{
			DisplayName: map[string]string{
				"en": "Event performed",
			},
		},
		Common: Common{
			Metadata: "{\"more\":\"data\"}",
		},
	}
)
