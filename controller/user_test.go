package controller

import "testing"

func TestPassword(t *testing.T) {
	password := "foobar"

	epw, err := passwordSecure(password)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := passwordCompare(password, epw)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := valid, true; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
