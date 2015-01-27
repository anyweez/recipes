package router

import (
	"errors"
	"fmt"
	"frontend/handlers"
	log "logging"
	"net/http"
)

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

func User(w http.ResponseWriter, r *http.Request) {
	route("/user", w, r, handlers.Registry{
		"GET": handlers.GetUser,
//		"POST": handlers.PostUser,
	})
}

/*
func UserLogin(w http.ResponseWriter, r *http.Request) {
	route("/user/login", w, r, handlers.Registry{
		"POST": handlers.PostLogUserIn,
	})
}

func Groups(w http.ResponseWriter, r *http.Request) {
	route("/groups", w, r, handlers.Registry{
		"GET": handlers.GetGroups,
	})
}
*/
