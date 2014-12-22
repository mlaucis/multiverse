/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

func getAccount(w http.ResponseWriter, r *http.Request) {
	var (
		accountId uint64
		err            error
	)
	vars := mux.Vars(r)

	if accountId, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		*entity.Account
	}{
		Account: &entity.Account{
			AccountID: accountId,
			Name:        "Demo Account",
			Enabled: true,
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
	}

	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getAccountUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountId     uint64
		userId string
		err       error
	)
	vars := mux.Vars(r)

	if accountId, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("userId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	userId = vars["userId"]

	response := &struct {
		*entity.AccountUser
	}{
		AccountUser: &entity.AccountUser{
			UserID:        userId,
			AccountID: accountId,
			Name: "Demo User",
			Email: "demouser@demo.com",
			Enabled: true,
			LastLogin:    "2014-12-20T12:10:10Z",
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
	}

	writeResponse(response, http.StatusOK, w)
}

func getAccountUserList(w http.ResponseWriter, r *http.Request) {
	var (
		accountId     uint64
		err       error
	)
	vars := mux.Vars(r)

	if accountId, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.Account
		AccountUser []*entity.AccountUser `json:"accountUser"`
	}{
		Account: entity.Account{
			AccountID: accountId,
			Name:        "Demo Account",
			Enabled: true,
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
		AccountUser: []*entity.AccountUser{
			&entity.AccountUser{
				UserID:        "1",
				AccountID: accountId,
				Name: "Demo User",
				Email: "demouser@demo.com",
				Enabled: true,
				LastLogin:    "2014-12-20T12:10:10Z",
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				UserID:        "2",
				AccountID: accountId,
				Name: "Demo User",
				Email: "demouser@demo.com",
				Enabled: true,
				LastLogin:    "2014-12-20T12:10:10Z",
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				UserID:        "3",
				AccountID: accountId,
				Name: "Demo User",
				Email: "demouser@demo.com",
				Enabled: true,
				LastLogin:    "2014-12-20T12:10:10Z",
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, w)
}

func getApplicationList(w http.ResponseWriter, r *http.Request) {
	var (
		accountId     uint64
		err       error
	)
	vars := mux.Vars(r)

	if accountId, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	response := &struct {
		entity.Account
		Application []*entity.Application `json:"application"`
	}{
		Account: entity.Account{
			AccountID: accountId,
			Name:        "Demo Account",
			Enabled: true,
			CreatedAt:    "2014-12-15T10:10:10Z",
			UpdatedAt:    "2014-12-20T12:10:10Z",
		},
		Application: []*entity.Application{
			&entity.Application{
				AppID:        1,
				Key: "abc123def",
				AccountID: accountId,
				Name: "Demo App",
				Enabled: true,
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				AppID:        2,
				Key: "abc345def",
				AccountID: accountId,
				Name: "Demo App",
				Enabled: true,
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				AppID:        3,
				Key: "abc678ef",
				AccountID: accountId,
				Name: "Demo App",
				Enabled: true,
				CreatedAt:    "2014-12-15T10:10:10Z",
				UpdatedAt:    "2014-12-20T12:10:10Z",
			},
		},
	}

	writeResponse(response, http.StatusOK, w)
}