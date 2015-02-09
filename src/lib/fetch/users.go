package fetch

import (
	gproto "code.google.com/p/goprotobuf/proto"	
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	proto "proto"
	"math/rand"
	"time"
)

var uc *mgo.Collection

func init() {
	conf, err := config.New("recipes.conf")
	if err != nil {
		log.Fatal("Couldn't load configuration file 'recipes.conf.'")
	}
	
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("User retrieval API can't connect to MongoDB instance: " + conf.Mongo.ConnectionString())
	}
	
	uc = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.UserCollection)
}

/**
 * Fetch a specific user by UserId.
 */
func UserById(user_id uint64) (proto.User, error) {
	var user proto.User
	err := uc.Find(bson.M{"id": user_id}).One(&user)
	
	return user, err
}

func UserByName(name string) (proto.User, error) {
	var user proto.User
	err := uc.Find(bson.M{"name": name}).One(&user)
	
	return user, err	
}

func CreateUser(user proto.User) (uint64, error) {
	rand.Seed( time.Now().UnixNano() )
	// TODO: replace this with something guaranteed to be unique
	user.Id = gproto.Uint64( uint64(rand.Int63()) )
	user.CreateMs = gproto.Int64( time.Now().Unix() * 1000 )
		
	err := uc.Insert(user)
	
	if err != nil {
		return 0, err
	}
	
	return *user.Id, nil
}
