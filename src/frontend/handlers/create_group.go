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

/**
 * Creates a new group and saves it to persistent storage. Groups are
 * denormalized in the backend and joined together during serving time.
 */
func CreateGroup(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
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

	// Get the group data from the post body
	grp := proto.Group{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&grp)

	// If there was an error decoding or a name is not provided, return with
	// an error indicating that there was an error with the data.
	if err != nil || len(*grp.Name) == 0 {
		le.Update(log.STATUS_ERROR, "Invalid post data: "+err.Error(), nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	// Get the user object
	gid, ferr := fetchme.CreateGroup(grp)

	if ferr != nil {
		le.Update(log.STATUS_ERROR, "Couldn't create group on persistent storage: "+ferr.Error(), nil)
	}

	le.Update(log.STATUS_OK, "Created group.", log.Fields{
		"gid": gid,
	})
}
