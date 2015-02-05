package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"fmt"
	"lib/fetch"
	log "logging"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
	// Check to make sure a user is logged in.
	if false {
		var user_id uint64 = 1
		// Get the current user
		user, err := fetch.UserById(user_id)
		
		if err != nil {
			le.Update(log.STATUS_ERROR, fmt.Sprintf("Unknown user ID: '%d'", user_id), nil)
			// TODO: add user-facing message
			return
		}
		
		// Encode as JSON
		data, _ := json.Marshal(user)
		w.Write( data )	
	} else {
		err := fee.WARNING_NOT_LOGGED_IN
		
		data, _ := json.Marshal(err)
		// Output
		w.WriteHeader(err.HttpCode)
		w.Write(data)
	}
}
