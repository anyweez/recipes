package handlers

import (
//	gproto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	fee "frontend/errors"
	"lib/fetch"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Check to make sure a user is logged in.
	if false {
		// Get the current user
		user := fetch.User(1)
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
