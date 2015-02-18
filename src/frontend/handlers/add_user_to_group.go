package handlers

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"fmt"
	fee "frontend/errors"
	"frontend/state"
	"github.com/gorilla/mux"

	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
	"strconv"
)

func AddUserToGroup(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
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

	// Get the user data from the post body
	user := proto.User{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	params := mux.Vars(r)
	groupid, cerr := strconv.ParseUint(params["group_id"], 10, 64)

	if cerr != nil {
		le.Update(log.STATUS_WARNING, "Invalid group ID.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	// If there was an error decoding or a name is not provided, return with
	// an error indicating that there was an error with the data.
	if err != nil || len(*user.EmailAddress) == 0 || groupid == 0 {
		fmt.Println(groupid)
		fmt.Println(user.EmailAddress)
		le.Update(log.STATUS_ERROR, "Invalid post data or missing parameter", nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	// Get the user object
	user, err = fetch.UserByEmail(*user.EmailAddress)

	if err != nil {
		le.Update(log.STATUS_ERROR, "User doesn't exist: "+err.Error(), nil)
		e := fee.USER_DOESNT_EXIST
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	ferr := fetch.AddUserToGroup(user, proto.Group{
		Id: gproto.Uint64(groupid),
	})

	if ferr != nil {
		le.Update(log.STATUS_ERROR, "Couldn't update group: "+ferr.Error(), nil)
		e := fee.COULDNT_COMPLETE_OPERATION
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	le.Update(log.STATUS_OK, "", log.Fields{
		"gid": groupid,
	})
}
