package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
)

type UserRequest struct {
	Name string
}

func GetUser(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
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
	session, serr := storage.Get(r, "userdata")

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values["user"]

	user, ferr := fetch.UserById(*ud.(*proto.User).Id)

	if ferr != nil {
		le.Update(log.STATUS_WARNING, "Couldn't read session data; the user doesn't seem to exist."+ferr.Error(), nil)
		err := fee.USER_DOESNT_EXIST
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	data, _ := json.Marshal(user)
	w.Write(data)
}
