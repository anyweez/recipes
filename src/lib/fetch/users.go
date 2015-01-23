package fetch

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	proto "proto"
)

var uc *mgo.Collection

func init() {
	conf := config.New("recipes.conf")
	
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("Ingredient retrieval API can't connect to MongoDB instance: " + conf.Mongo.ConnectionString())
	}
	
	uc = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.UserCollection)
}

/**
 * Fetch a specific user by UserId.
 */
func User(user_id uint64) proto.User {
	var user proto.User
	uc.Find(bson.M{"id": user_id}).One(&user)
	
	return user
}
