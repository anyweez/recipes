package fetch

import (
	//	gproto "code.google.com/p/goprotobuf/proto
	"lib/config"
	"labix.org/v2/mgo"
	"log"
	proto "proto"
)

var ic *mgo.Collection

func init() {
	conf, _ := config.New("recipes.conf")
	
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("Ingredient retrieval API can't connect to MongoDB instance: " + conf.Mongo.ConnectionString())
	}
	
	ic = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.IngredientCollection)
}

/**
 * This function retrieves a list of ingredients and is designed to be
 * the canonical API for the underlying data storage layer.
 */
func AllIngredients() []proto.Ingredient {
	ingredients := make([]proto.Ingredient, 0)
	iter := ic.Find(nil).Iter()
	
	result := proto.Ingredient{}
	for iter.Next(&result) {
		ingredients = append(ingredients, result)
	}

	return ingredients
}
