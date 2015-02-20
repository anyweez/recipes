package handlers

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	fee "frontend/errors"
	"frontend/state"
	"lib/fetch"
	log "logging"
	"math/rand"
	"net/http"
	proto "proto"
	"time"
)

type RecipeVoteRequest struct {
	Vote   bool
	Recipe string
	Group  uint64
}

func SetMealVote(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
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

	// Get parameters from the post body
	rvr := RecipeVoteRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rvr)

	if err != nil {
		le.Update(log.STATUS_ERROR, "Invalid post data: "+err.Error(), nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}

	status := proto.RecipeVote_YES
	if !rvr.Vote {
		status = proto.RecipeVote_NO
	}

	session, serr := ss.Session.Get(r, "userdata")

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values[state.UserDataActiveUser]
	user, _ := fetchme.UserById(*ud.(*proto.User).Id)

	rand.Seed(time.Now().Unix())

	group := proto.Group{
		Id: gproto.Uint64(rvr.Group),
	}

	// TODO: handle this error.
	meal, _ := fetchme.GetCurrentMeal(group)
	recipe := fetchme.Recipe(rvr.Recipe)
	nu := fetchme.NormalizeUser(user)

	vote := proto.RecipeVote{
		Id:     gproto.Uint64(uint64(rand.Uint32())),
		User:   &nu,
		Group:  &group,
		Meal:   &meal,
		Recipe: &recipe,
		Status: &status,
	}

	// Store the vote.
	fetchme.StoreVote(vote)

	// Check to see whether there's agreement on this recipe.
	// If so, copy the recipe into the meal object.
	if fetchme.CheckForQuorum(vote) {
		meal.Recipe = &recipe

		fetchme.UpdateMeal(meal)
	}
}
