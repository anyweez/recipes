package fetch

import (
	//	gproto "code.google.com/p/goprotobuf/proto"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	proto "proto"
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

	uc = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.GroupCollection)
}

func GroupsForUser(u proto.User) ([]*proto.Group, error) {
	var groups []*proto.Group

	err := uc.Find(bson.M{"users": bson.M{"id": *u.Id}}).All(&groups)

	return groups, err
}
