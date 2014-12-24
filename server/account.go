/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

func getAccount(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		*entity.Account
	}{
		Account: &entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getAccountUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		userID    string
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userID = vars["userId"]

	response := &struct {
		*entity.AccountUser
	}{
		AccountUser: &entity.AccountUser{
			ID:        userID,
			AccountID: accountID,
			Name:      "Demo User",
			Email:     "demouser@demo.com",
			Enabled:   true,
			LastLogin: "2014-12-20T12:10:10Z",
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getAccountUserList(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.Account
		AccountUser []*entity.AccountUser `json:"accountUser"`
	}{
		Account: entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
		AccountUser: []*entity.AccountUser{
			&entity.AccountUser{
				ID:        "1",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				ID:        "2",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				ID:        "3",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}

func getAccountApplications(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.Account
		Application []*entity.Application `json:"application"`
	}{
		Account: entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
		Application: []*entity.Application{
			&entity.Application{
				ID:        1,
				Key:       "abc123def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				ID:        2,
				Key:       "abc345def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				ID:        3,
				Key:       "abc678ef",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, 10, w, r)
}
