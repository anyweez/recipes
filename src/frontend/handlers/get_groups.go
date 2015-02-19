package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"frontend/state"
	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
)

func GetGroups(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
	fetchme := fetch.NewFetcher(ss)

	// If the requested user isn't logged in there's nothing we can do
	// for them.
	if !IsLoggedIn(ss, r) {
		le.Update(log.STATUS_WARNING, "User not logged in.", nil)
		err := fee.NOT_LOGGED_IN
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	// Retrieve the session
	session, serr := ss.Session.Get(r, state.UserDataSession)

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values[state.UserDataActiveUser]
	groups, ferr := fetchme.GroupsForUser(*ud.(*proto.User))

	if ferr != nil {
		le.Update(log.STATUS_ERROR, "Couldn't retrieve groups from database: "+ferr.Error(), nil)
	}

	data, _ := json.Marshal(groups)
	w.Write(data)
}
