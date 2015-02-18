package handlers

import (
	"encoding/json"
	"fmt"
	fee "frontend/errors"
	"frontend/state"
	"lib/fetch"
	log "logging"
	"net/http"
)

type LoginRequest struct {
	Name         string
	EmailAddress string
}

func Login(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
	post_request := LoginRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post_request)

	// If we can't read the body, throw an error.
	if err != nil {
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)
		w.WriteHeader(e.HttpCode)
		w.Write(data)

		return
	}

	session, serr := storage.Get(r, "userdata")

	// If the session couldn't be decoded, we've got to return an error.
	// This shouldn't happen unless something were to go wrong.
	if serr != nil {
		//	log.Println( serr.Error() )
		cserr := fee.CORRUPTED_SESSION
		data, _ := json.Marshal(cserr)

		w.WriteHeader(cserr.HttpCode)
		w.Write(data)
	}

	// Check if the user is logged in already.
	if _, exists := session.Values[UserDataActiveUser].(string); exists {
		le.Update(log.STATUS_WARNING, fmt.Sprintf("User is already logged in.", post_request.EmailAddress), nil)
		return
	}

	// Store the user's data in the encrypted session.
	// TODO: validate the email address.
	user, err := fetch.UserByEmail(post_request.EmailAddress)

	// If the user doesn't exist, return an error.
	if err != nil {
		le.Update(log.STATUS_ERROR, "The requested user couldn't be found: "+err.Error(), nil)

		e := fee.USER_DOESNT_EXIST
		data, _ := json.Marshal(e)
		w.WriteHeader(e.HttpCode)
		w.Write(data)
		// If the user does exist, store their data in the session.
	} else {
		session.Values[UserDataActiveUser] = user
		werr := session.Save(r, w)

		if werr != nil {
			le.Update(log.STATUS_ERROR, "Error storing user data in session: "+werr.Error(), nil)
		}
	}
}
