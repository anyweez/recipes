package fetch

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"labix.org/v2/mgo/bson"
	"math/rand"
	proto "proto"
	"time"
)

/**
 * Fetch a specific user by UserId.
 */
func (f *Fetcher) UserById(user_id uint64) (proto.User, error) {
	var user proto.User
	err := f.SS.Database.Users.Find(bson.M{"id": user_id}).One(&user)

	return user, err
}

// TODO: Hoping this can be a temporary function.
func (f *Fetcher) UserByEmail(name string) (proto.User, error) {
	var user proto.User
	err := f.SS.Database.Users.Find(bson.M{"emailaddress": name}).One(&user)

	return user, err
}

func (f *Fetcher) CreateUser(user proto.User) (uint64, error) {
	rand.Seed(time.Now().UnixNano())
	// TODO: replace this with something guaranteed to be unique
	user.Id = gproto.Uint64(uint64(rand.Int63()))
	user.CreateMs = gproto.Int64(time.Now().Unix() * 1000)

	err := f.SS.Database.Users.Insert(user)

	if err != nil {
		return 0, err
	}

	return *user.Id, nil
}

func (f *Fetcher) NormalizeUser(user proto.User) proto.User {
	user.Name = gproto.String("")
	user.EmailAddress = gproto.String("")
	user.CreateMs = gproto.Int64(0)

	return user
}
