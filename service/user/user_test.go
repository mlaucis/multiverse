package user

import "testing"

const (
	validEmail    = "user123@tp.gl"
	validPassword = "1234"
)

func TestValidate(t *testing.T) {
	us := List{
		{},                                                                          // Email and Username missing
		{Email: "user@foo"},                                                         // Email invalid
		{Email: validEmail, Firstname: randStringRunes(1)},                          // Firstname min length
		{Email: validEmail, Firstname: randStringRunes(41)},                         // Firstname max length
		{Email: validEmail, Lastname: randStringRunes(1)},                           // Lastname min length
		{Email: validEmail, Lastname: randStringRunes(41)},                          // Lastname max length
		{Email: validEmail, Password: validPassword, Username: randStringRunes(1)},  // Username min length
		{Email: validEmail, Password: validPassword, Username: randStringRunes(41)}, // Username max length
		{Email: validEmail, Password: ""},                                           // Password empty
		{Email: validEmail, URL: "foo\bar"},                                         // URL invalid
	}

	for _, u := range us {
		if have, want := u.Validate(), ErrInvalidUser; !IsInvalidUser(have) {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
