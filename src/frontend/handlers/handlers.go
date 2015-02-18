package handlers

import (
	"fmt"
	"frontend/state"
	"net/http"
)

/**
 * A handler is a function that performs an action based on a GET or
 * POST request and returns the status of the operation to the frontend.
 *
 * It is a generic interface for other functions in this package.
 */
type Handler func(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent)

/**
 * A function to determine whether a user with a given name is logged in.
 */
func IsLoggedIn(r *http.Request) bool {
	session, serr := storage.Get(r, UserDataSession)

	if serr != nil {
		fmt.Println("IsLoggedIn: " + serr.Error())
	}

	_, exists := session.Values[UserDataActiveUser]
	return exists
}
