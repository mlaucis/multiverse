package controller

import (
	"fmt"
	"testing"

	"github.com/tapglue/multiverse/service/user"
)

func TestFilterUsers(t *testing.T) {
	us := testUsers(10)

	us, err := filterUsers(us, testConditionUserEven)
	if err != nil {
		t.Fatal(err)
	}

	if have, want := len(us), 5; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

func TestFilterUsersError(t *testing.T) {
	_, err := filterUsers(testUsers(1), testConditionUserError)
	if err == nil {
		t.Error("want error")
	}
}

func testConditionUserEven(user *user.User) (bool, error) {
	return user.ID%2 == 0, nil
}

func testConditionUserError(user *user.User) (bool, error) {
	return false, fmt.Errorf("condition errored")
}

func testUsers(n int) user.List {
	us := user.List{}

	for i := 0; i < n; i++ {
		us = append(us, &user.User{
			ID: uint64(i + 1),
		})
	}

	return us
}
