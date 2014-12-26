/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gluee/backend/db"
	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

/**
 * getAccount handles requests to a single account
 * Request: GET /account/:AccountID
 * Test with: curl -i localhost/account/:AccountID
 * @param w, response writer
 * @param r, http request
 */
func getAccount(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		account   = &entity.Account{}
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Read account from database
	err = db.GetSlave().QueryRowx("SELECT * FROM accounts WHERE id=?", accountID).StructScan(account)
	if err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusOK, 10, w, r)
}

/**
 * createAccount handles requests create an account
 * Request: POST /account
 * Test with: curl -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/account
 * @param w, response writer
 * @param r, http request
 */
func createAccount(w http.ResponseWriter, r *http.Request) {
	//INSERT INTO `gluee`.`accounts` (`id`, `name`, `enabled`, `created_at`, `updated_at`) VALUES (NULL, 'demo', '1', '2014-12-26 11:14:24', '2014-12-26 11:14:24');
}
