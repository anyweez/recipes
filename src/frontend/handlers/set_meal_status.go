package handlers

import (
	"encoding/json"
	gproto "code.google.com/p/goprotobuf/proto"
	fee "frontend/errors"
	log "logging"
	"lib/fetch"
	"net/http"
	proto "proto"
)

type MealStatus struct {
	Group		uint64
	Available	bool
}

func SetMealStatus(w http.ResponseWriter, r *http.Request, le log.LogEvent) {
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
	
	// Get parameters from the post body
	ms := MealStatus{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ms)

	if err != nil {
		le.Update(log.STATUS_ERROR, "Invalid post data: "+err.Error(), nil)
		e := fee.INVALID_POST_DATA
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}
	
	status := proto.RecipeVote_NOT_SET
	if !ms.Available {
		status = proto.RecipeVote_ABSTAIN
	}
	
	meal, merr := fetch.GetCurrentMeal(proto.Group{
		Id: gproto.Uint64(ms.Group),
	})
	
	if merr != nil {
		le.Update(log.STATUS_ERROR, "Invalid post data: "+merr.Error(), nil)
		e := fee.COULDNT_COMPLETE_OPERATION
		data, _ := json.Marshal(e)

		w.WriteHeader(e.HttpCode)
		w.Write(data)
		return
	}
	
	session, serr := storage.Get(r, "userdata")

	if serr != nil {
		le.Update(log.STATUS_WARNING, "User data doesn't exist for logged in user:"+serr.Error(), nil)
		return
	}

	// Get the user object
	ud, _ := session.Values[UserDataActiveUser]
	user, _ := fetch.UserById(*ud.(*proto.User).Id)
	
	// Check to see if a vote already exists. If so, set it to ABSTAIN.
	// If not, create a new vote object and mark it as ABSTAIN.
	added := false
	for i := 0; i < len(meal.Votes); i++ {
		if *user.Id == *meal.Votes[i].User.Id {
			meal.Votes[i].Status = &status
			added = true
		}
	}
	
	if !added {
		user := fetch.NormalizeUser(user)
		group := fetch.NormalizeGroup(proto.Group{
				Id: gproto.Uint64(ms.Group),
		})
		
		meal.Votes = append(meal.Votes, &proto.RecipeVote{
			User: &user,
			Group: &group,
			Status: &status,
		})
	}
	
	fetch.UpdateMeal(meal)
}
