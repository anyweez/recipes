package handlers

import (
	"fmt"
	"frontend/state"
	log "logging"
	"net/http"
	proto "proto"
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
func IsLoggedIn(ss *state.SharedState, r *http.Request) bool {
	session, serr := ss.Session.Get(r, state.UserDataSession)

	if serr != nil {
		fmt.Println("IsLoggedIn: " + serr.Error())
	}

	_, exists := session.Values[state.UserDataActiveUser]
	return exists
}

/**
 * Retrieves information about the logged in user from the session
 * if a user is logged in.
 */
func GetLoggedInUser(ss *state.SharedState, r *http.Request) (proto.User, bool) {
	// If the user is logged in, fetch and return the user object.
	if IsLoggedIn(ss, r) {
		session, err := ss.Session.Get(r, state.UserDataSession)

		if err != nil {
			return proto.User{}, false
		}

		user, exists := session.Values[state.UserDataActiveUser]
		if exists {
			return user.(proto.User), true
		}
	}
	// If the user is not logged on, return an empty User object.
	return proto.User{}, false
}
