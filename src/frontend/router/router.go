package router

import (
	"errors"
	"fmt"
	"frontend/handlers"
	log "logging"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request) {
	route("/users", w, r, handlers.Registry{
		"POST": handlers.CreateNewUser,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	route("/users/login", w, r, handlers.Registry{
		"POST": handlers.Login,
	})
}

func UserMe(w http.ResponseWriter, r *http.Request) {
	route("/users/me", w, r, handlers.Registry{
		"GET": handlers.GetUser,
	})
}

func Groups(w http.ResponseWriter, r *http.Request) {
	route("/groups", w, r, handlers.Registry{
		"GET": handlers.GetGroups,
		"POST": handlers.CreateGroup,
	})
}

func GroupsIdUser(w http.ResponseWriter, r *http.Request) {
	route("/groups/{group_id}/user", w, r, handlers.Registry{
		"POST": handlers.AddUserToGroup,
	})
}

func MealsToday(w http.ResponseWriter, r *http.Request) {
	route("/meals/today", w, r, handlers.Registry{
		"GET": handlers.GetTodaysMeal,
	})
}

func MealsTodayStatus(w http.ResponseWriter, r *http.Request) {
	route("/meals/today/status", w, r, handlers.Registry{
		"POST": handlers.SetMealStatus,
	})
}

func MealsVote(w http.ResponseWriter, r *http.Request) {
	route("/meals/vote", w, r, handlers.Registry{
		"POST": handlers.SetMealVote,
	})
}

func Recipes(w http.ResponseWriter, r *http.Request) {
	route("/recipes", w, r, handlers.Registry{
		"GET": handlers.GetBestRecipes,
	})
}

/**
 * General-purpose function that handlers common routing code between all individual routing functions.
 * All routing functions should call this function with inputs specific to their user case (and should
 * not handle the routing themselves.
 */
func route(path string, w http.ResponseWriter, r *http.Request, hndl handlers.Registry) error {
	le := log.New("web_request", log.Fields{
		"handled_path": path,
		"method": r.Method,
	})
	fn, exists := hndl[r.Method]
	
	if exists {
		fn(w, r)
		le.Update(log.STATUS_COMPLETE, "", nil)
		
		return nil
	} else {
		msg := fmt.Sprintf("No handler specified for method %s on path %s", r.Method, path)
		
		le.Update(log.STATUS_ERROR, msg, nil)
		return errors.New( msg )
	}
}
