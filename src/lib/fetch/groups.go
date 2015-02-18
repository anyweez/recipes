package fetch

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
	"math/rand"
	proto "proto"
	"time"
)

func (f *Fetcher) GroupById(group_id uint64) (proto.Group, error) {
	var group proto.Group
	// TODO: change this to a real number.
	err := f.SS.Database.Groups.Find(bson.M{"id": group_id}).One(&group)

	if err != nil {
		fmt.Println(fmt.Sprintf("Error retrieving group #%d", group_id))
		return proto.Group{}, err
	}

	return f.DenormalizeUsers(group), err
}

func (f *Fetcher) GroupsForUser(u proto.User) ([]*proto.Group, error) {
	var groups []*proto.Group

	err := f.SS.Database.Groups.Find(bson.M{"members": bson.M{"$elemMatch": bson.M{"id": *u.Id}}}).All(&groups)

	// Fetch all user data for the groups.

	for gi := 0; gi < len(groups); gi++ {
		group := f.DenormalizeUsers(*groups[gi])
		groups[gi] = &group
	}

	return groups, err
}

func (f *Fetcher) CreateGroup(g proto.Group) (uint64, error) {
	rand.Seed(time.Now().UnixNano())
	// TODO: replace this with something guaranteed to be unique
	g.Id = gproto.Uint64(uint64(rand.Int63()))
	g.CreateMs = gproto.Int64(time.Now().Unix() * 1000)

	// Clear out the fields that shouldn't be stored.
	for i := 0; i < len(g.Members); i++ {
		nu := f.NormalizeUser(*g.Members[i])
		g.Members[i] = &nu
	}

	err := f.SS.Database.Groups.Insert(g)

	if err != nil {
		return 0, err
	}

	return *g.Id, nil
}

func (f *Fetcher) AddUserToGroup(u proto.User, g proto.Group) error {
	nu := f.NormalizeUser(u)

	// Fetch the group
	group, err := f.GroupById(*g.Id)
	if err != nil {
		return err
	}

	exists := false
	for i := 0; i < len(group.Members); i++ {
		if *group.Members[i].Id == *u.Id {
			exists = true
		}
	}

	if !exists {
		// Add the normalized user to the group.
		group.Members = append(group.Members, &nu)

		return f.SS.Database.Groups.Update(bson.M{"id": group.Id}, group)
	} else {
		return errors.New("User already a member of group.")
	}
}

func (f *Fetcher) NormalizeGroup(g proto.Group) proto.Group {
	g.Name = gproto.String("")
	g.CreateMs = gproto.Int64(0)

	for i := 0; i < len(g.Members); i++ {
		user := f.NormalizeUser(*g.Members[i])
		g.Members[i] = &user
	}

	return g
}

/**
 * Fetch user information for each user in the specified group.
 */
func (f *Fetcher) DenormalizeUsers(group proto.Group) proto.Group {
	for ui := 0; ui < len(group.Members); ui++ {
		user, _ := f.UserById(*group.Members[ui].Id)
		group.Members[ui] = &user
	}

	return group
}
