package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
)

func GetGroups(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
	// If the requested user isn't logged in there's nothing we can do
	// for them.
	if !IsLoggedIn(r) {
		le.Update(log.STATUS_WARNING, "User not logged in.", nil)
		err := fee.NOT_LOGGED_IN
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	// Retrieve the session
	session, serr := storage.Get(r, UserDataSession)

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values[UserDataActiveUser]
	groups, ferr := fetch.GroupsForUser(*ud.(*proto.User))

	if ferr != nil {
		le.Update(log.STATUS_ERROR, "Couldn't retrieve groups from database: "+ferr.Error(), nil)
	}

	data, _ := json.Marshal(groups)
	w.Write(data)
}
