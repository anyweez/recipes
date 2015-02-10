package fetch

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	"math/rand"
	proto "proto"
	"time"
)

var gc *mgo.Collection

func init() {
	conf, err := config.New("recipes.conf")
	if err != nil {
		log.Fatal("Couldn't load configuration file 'recipes.conf.'")
	}

	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("User retrieval API can't connect to MongoDB instance: " + conf.Mongo.ConnectionString())
	}

	gc = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.GroupCollection)
}

func GroupsForUser(u proto.User) ([]*proto.Group, error) {
	var groups []*proto.Group

	err := gc.Find(bson.M{"members": bson.M{"$elemMatch": bson.M{"id": *u.Id}}}).All(&groups)

	// Fetch all user data for the groups.
	for gi := 0; gi < len(groups); gi++ {
		for ui := 0; ui < len(groups[gi].Members); ui++ {
			user, _ := UserById(*groups[gi].Members[ui].Id)
			groups[gi].Members[ui] = &user
		}
	}

	return groups, err
}

func CreateGroup(g proto.Group) (uint64, error) {
	rand.Seed(time.Now().UnixNano())
	// TODO: replace this with something guaranteed to be unique
	g.Id = gproto.Uint64(uint64(rand.Int63()))
	g.CreateMs = gproto.Int64(time.Now().Unix() * 1000)

	// Clear out the fields that shouldn't be stored.
	for i := 0; i < len(g.Members); i++ {
		nu := NormalizeUser(*g.Members[i])
		g.Members[i] = &nu
	}

	err := gc.Insert(g)

	if err != nil {
		return 0, err
	}

	return *g.Id, nil
}
