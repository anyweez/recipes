package handlers

import (
"fmt"
	"github.com/gorilla/sessions"
//	"github.com/gorilla/securecookie"
	log "logging"
	"net/http"
)

var storage = sessions.NewCookieStore(
	[]byte("hello"),
//	securecookie.GenerateRandomKey(32),	// Authentication
//	securecookie.GenerateRandomKey(32),	// Encryption
)

func init() {
	storage.Options = &sessions.Options{
//		Domain: "localhost",
//		Path: "/",
		MaxAge: 3600 * 365,	// 1 year
		HttpOnly: true,
	}
}

/**
 * A handler is a function that performs an action based on a GET or
 * POST request and returns the status of the operation to the frontend.
 * 
 * It is a generic interface for other functions in this package.
 */
type Handler func(w http.ResponseWriter, r *http.Request, le log.LogEvent) 

/**
 * Registries contain mappings between HTTP methods (GET, POST, etc) and
 * the handlers that should be used to fulfill the request.
 */
type Registry map[string]Handler

//var store = sessions.NewCookieStore()

/**
 * A function to determine whether a user with a given name is logged in.
 */
func IsLoggedIn(r *http.Request) bool {
	session, serr := storage.Get(r, "userdata")
		
	if serr != nil {
		fmt.Println("IsLoggedIn: " + serr.Error())
	}
		
	_, exists := session.Values["user"]
	return exists
}
