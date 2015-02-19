package state

import (
	"encoding/gob"
	"github.com/gorilla/sessions"
	"labix.org/v2/mgo"
	"lib/config"
	log "logging"
	"net/rpc"
	proto "proto"
)

const (
	// Session key for Session elements
	UserDataSession = "userdata"

	// Fields within a session element
	UserDataActiveUser = "user"
)

type SharedState struct {
	// The configuration object, specified by the user, that contains a ton
	// of parameters that are used for reading / writing data.
	Config *config.RecipesConfig
	// Encrypted user session object.
	Session *sessions.CookieStore
	// A bunch of database connections. Note that only *shared* database
	// connections should be stored here; if a connection shouldn't be shared
	// then it shouldn't be stored here.
	Database *DBState
	// RPC client used to distribute computationally intensive tasks to backend
	// nodes.
	Retriever *rpc.Client
}

type DBState struct {
	Groups      *mgo.Collection
	Ingredients *mgo.Collection
	Meals       *mgo.Collection
	Recipes     *mgo.Collection
	Users       *mgo.Collection
	Votes       *mgo.Collection
}

/**
 * The function initializes all shared state that's required for handling requests made to
 * the frontend (database connections, configuration object, etc).
 */
func NewSharedState(conf config.RecipesConfig) (*SharedState, error) {
	le := log.New("shared_state", nil)
	ss := SharedState{}
	var err error

	// Set the configuration object.
	ss.Config = &conf
	ss.Session = createSession()
	ss.Database, err = connectToDB(conf)

	// If we can't connect to the database, fail.
	if err != nil {
		le.Update(log.STATUS_ERROR, "Couldn't connect to database: "+err.Error(), nil)
		return &ss, err
	}

	ss.Retriever, err = rpc.DialHTTP("tcp", conf.Rpc.ConnectionString())

	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't connect to retriever: "+err.Error(), nil)
		return &ss, err
	}

	return &ss, nil
}

func createSession() *sessions.CookieStore {
	storage := sessions.NewCookieStore(
		[]byte("hello"),
	//	securecookie.GenerateRandomKey(32),	// Authentication
	//	securecookie.GenerateRandomKey(32),	// Encryption
	)

	storage.Options = &sessions.Options{
		//		Domain: "localhost",
		Path:     "/",
		MaxAge:   3600 * 365, // 1 year
		HttpOnly: true,
	}

	// Register users to be encodable as gobs so that they can be stored
	// in sessions.
	gob.Register(&proto.User{})

	return storage
}

func connectToDB(conf config.RecipesConfig) (*DBState, error) {
	dbs := DBState{}

	// TODO: check to make sure that we successfully connected to the specified database
	session, err := mgo.Dial(conf.Mongo.ConnectionString())

	if err != nil {
		return &dbs, err
	}

	dbs.Groups = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.GroupCollection)
	dbs.Ingredients = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.IngredientCollection)
	dbs.Meals = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.MealsCollection)
	dbs.Recipes = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
	dbs.Users = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.UserCollection)
	dbs.Votes = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.VotesCollection)

	return &dbs, nil
}
