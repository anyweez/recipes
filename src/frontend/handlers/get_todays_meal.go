package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	gproto "code.google.com/p/goprotobuf/proto"
	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
	"strconv"
)

func GetTodaysMeal(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
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
	
	gid, exists := r.URL.Query()["group"]
	
	if !exists {
		le.Update(log.STATUS_WARNING, "No group ID specified.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}
	
	// Parse the string into a number.
	groupid, perr := strconv.ParseUint(gid[0], 10, 64)
	
	if perr != nil {
		le.Update(log.STATUS_WARNING, "Invalid group ID:"+perr.Error(), nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}
	
	meal, err := fetch.GetCurrentMeal(proto.Group{
		Id: gproto.Uint64(groupid),
	})
	
	if err != nil {
		le.Update(log.STATUS_WARNING, "Couldn't fetch current meal:"+err.Error(), nil)
		err := fee.COULDNT_COMPLETE_OPERATION
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	data, _ := json.Marshal(meal)
	w.Write(data)
}
