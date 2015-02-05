package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	fetch "lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
)

/**
 * This method creates a new user with the information included in the
 * body of the request.
 */

func CreateNewUser(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
	// Get the user from the post body
	user := proto.User{}
		
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		le.Update(log.STATUS_ERROR, "Invalid post data: " + err.Error(), nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)
		
		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}
	
	// TODO: check to make sure it only includes valid characters
	_, err = fetch.CreateUser(user)
	
	if err != nil {
		le.Update(log.STATUS_ERROR, "Couldn't create user.", nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)
		
		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(""))
}

