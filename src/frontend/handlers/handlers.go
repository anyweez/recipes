package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	//	"github.com/gorilla/securecookie"
	"lib/config"
	log "logging"
	"net/http"
	"net/rpc"
	proto "proto"
)

var storage = sessions.NewCookieStore(
	[]byte("hello"),

//	securecookie.GenerateRandomKey(32),	// Authentication
//	securecookie.GenerateRandomKey(32),	// Encryption
)

// The recipe client (RPC).
var res *rpc.Client

const (
	// Session elements
	UserDataSession = "userdata"

	// Fields within a session element
	UserDataActiveUser = "user"
)

func init() {
	le := log.New("init_handlers", nil)

	// TODO: need to replace this.
	conf, err := config.New("recipes.conf")
	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't read configuration.", nil)
		return
	}
	
	storage.Options = &sessions.Options{
		//		Domain: "localhost",
		Path:     "/",
		MaxAge:   3600 * 365, // 1 year
		HttpOnly: true,
	}

	// Register users to be encodable as gobs so that they can be stored
	// in sessions.
	gob.Register(&proto.User{})
	
	res, err = rpc.DialHTTP("tcp", conf.Rpc.ConnectionString())
	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't connect to retriever: "+err.Error(), nil)
		return
	}

	le.Update(log.STATUS_COMPLETE, "", nil)
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

/**
 * A function to determine whether a user with a given name is logged in.
 */
func IsLoggedIn(r *http.Request) bool {
	session, serr := storage.Get(r, UserDataSession)

	if serr != nil {
		fmt.Println("IsLoggedIn: " + serr.Error())
	}

	_, exists := session.Values[UserDataActiveUser]
	return exists
}
