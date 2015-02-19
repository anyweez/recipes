package fetch

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"errors"
	"fmt"
	//	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"lib/config"
	//	"log"
	"math/rand"
	proto "proto"
	"strconv"
	"time"
)

/**
 * Gets the meal object for today, unless it doesn't exist, in which case
 * it creates a new object and stores it before returning it. Note that
 * the current design limits the system to one meal per day per group.
 */
func (f *Fetcher) GetCurrentMeal(g proto.Group) (proto.Meal, error) {
	var meal proto.Meal
	// Generate the datecode: YYYYMMDD
	dc, _ := strconv.Atoi(time.Now().Format("20060102"))
	datecode := int32(dc)

	fmt.Println(datecode)
	fmt.Println(*g.Id)
	err := f.SS.Database.Meals.Find(bson.M{"datecode": datecode, "group.id": *g.Id}).One(&meal)

	// If meal doesn't exist, create a new one.
	if err != nil {
		group := f.NormalizeGroup(g)

		rand.Seed(time.Now().UnixNano())
		meal = proto.Meal{
			// TODO: replace this with something guaranteed to be unique
			Id:       gproto.Uint64(uint64(rand.Int63())),
			Group:    &group,
			Datecode: gproto.Int32(datecode),
		}

		// Store this for next time.
		f.SS.Database.Meals.Insert(meal)
	}

	grp, gerr := f.GroupById(*g.Id)
	meal.Group = &grp

	if gerr != nil {
		return proto.Meal{}, errors.New("Group no longer exists:" + gerr.Error())
	}

	return meal, nil
}

func (f *Fetcher) UpdateMeal(m proto.Meal) {
	fmt.Println(*m.Id)
	grp := f.NormalizeGroup(*m.Group)
	m.Group = &grp
	f.SS.Database.Meals.Update(bson.M{"id": *m.Id}, m)
}
