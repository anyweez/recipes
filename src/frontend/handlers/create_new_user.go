package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"net/http"
	proto "proto"
)

/**
 * This method creates a new user with the information included in the
 * body of the request.
 */

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
	// Get the user from the post body
	user := proto.User{}
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)
		
		w.WriteHeader(e.HttpCode)
		w.Write(data)
	}
	
	// TODO: check to make sure it only includes valid characters
	
	// Re-encode into JSON to insert into storage.
	data, _ := json.Marshal(user)
	w.Write( data )
}

