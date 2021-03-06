package handlers

import (
	"encoding/json"
	fee "frontend/errors"
	"frontend/state"
	log "logging"
	"math/rand"
	"net/http"
	proto "proto"
	retrieve "retrieve"
	"strconv"
	"time"
)

type RecipeRequest struct {
	GroupId uint64
	Count   int
}

// TODO: clean up error handling here. There must be a better way once patterns emerge.
func GetBestRecipes(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
	if !IsLoggedIn(ss, r) {
		le.Update(log.STATUS_WARNING, "User not logged in.", nil)
		err := fee.NOT_LOGGED_IN
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	qry := r.URL.Query()

	gp, exists := qry["group"]
	// If the param doesn't exist, error.
	if len(gp) == 0 {
		le.Update(log.STATUS_ERROR, "Invalid fields provided in get request.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	// Check for the existence and proper type of the group id
	if _, terr := strconv.Atoi(gp[0]); terr != nil || !exists {
		le.Update(log.STATUS_ERROR, "Invalid fields provided in get request.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	tgid, _ := strconv.Atoi(gp[0])
	groupid := uint64(tgid)

	cp, exists := qry["count"]
	// If the param doesn't exist, error.
	if len(cp) == 0 {
		le.Update(log.STATUS_ERROR, "Invalid fields provided in get request.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}

	// Check for the existence and proper type of count.
	if _, terr := strconv.Atoi(cp[0]); terr != nil {
		le.Update(log.STATUS_ERROR, "Invalid fields provided in get request.", nil)
		err := fee.MISSING_QUERY_PARAMS
		data, _ := json.Marshal(err)

		w.WriteHeader(err.HttpCode)
		w.Write(data)
		return
	}
	count, _ := strconv.Atoi(cp[0])

	// Retrieve the session
	session, serr := ss.Session.Get(r, state.UserDataSession)

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values[state.UserDataActiveUser]
	// Generate a random seed used to specify which recipes should be
	// selected lacking stronger signals.
	// TODO: move this to serverside.
	rand.Seed(time.Now().UnixNano())
	seed := int64(rand.Int())

	recipes := make([]proto.Recipe, 0)
	err := ss.Retriever.Call("Retriever.GetBestRecipes", retrieve.BestRecipesRequest{
		Seed:    seed,
		UserId:  *ud.(*proto.User).Id,
		GroupId: groupid,
		Count:   count,
	}, &recipes)

	if err != nil {
		le.Update(log.STATUS_ERROR, "Invalid or no response from RPC call.", nil)
	}

	data, _ := json.Marshal(recipes)
	w.Write(data)

	le.Update(log.STATUS_COMPLETE, "", nil)
}
