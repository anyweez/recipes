package ingredients

import (
	//	gproto "code.google.com/p/goprotobuf/proto
	"labix.org/v2/mgo"
	"log"
	proto "proto"
)

/**
 * This function retrieves a list of ingredients and is designed to be
 * the canonical API for the underlying data storage layer.
 */
func GetAll() []proto.Ingredient {
	log.Println("Loading ingredient list...")

	// TODO: Need to base this off of a variable.
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB instance.")
	}
	c := session.DB("recipes").C("ingredients")

	iter := c.Find(nil).Iter()
	result := proto.Ingredient{}

	ingredients := make([]proto.Ingredient, 0)
	for iter.Next(&result) {
		ingredients = append(ingredients, result)
	}

	return ingredients
}
