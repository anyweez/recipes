package handlers

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	fee "frontend/errors"
	"frontend/state"
	"lib/fetch"
	log "logging"
	"net/http"
	proto "proto"
	"strconv"
)

func GetTodaysMeal(w http.ResponseWriter, r *http.Request, ss *state.SharedState, le log.LogEvent) {
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

	meal, err := fetchme.GetCurrentMeal(proto.Group{
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

	votes := fetchme.VotesForMeal(meal)
	// For each member of the group, check to see if they've abstained.
	// If not, see if a vote's been cast yet.
	for i := 0; i < len(meal.Group.Members); i++ {
		// Check to see if the user has abstained for this meal.
		abstained := false
		for k := 0; k < len(meal.Votes); k++ {
			if *meal.Votes[k].User.Id == *meal.Group.Members[i].Id && *meal.Votes[k].Status == proto.RecipeVote_ABSTAIN {
				abstained = true
			}
		}

		// If not, check to see if their vote has been cast.
		if !abstained {
			for j := 0; j < len(votes); j++ {
				for k := 0; k < len(votes); k++ {
					if *votes[k].User.Id == *meal.Group.Members[i].Id {
						meal.Votes = append(meal.Votes, &votes[k])
					}
				}
			}
		}
	}

	data, _ := json.Marshal(meal)
	w.Write(data)
}
