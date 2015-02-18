package fetch

import (
	//	gproto "code.google.com/p/goprotobuf/proto
	//	"labix.org/v2/mgo"
	//	"lib/config"
	//	"log"
	proto "proto"
)

/**
 * This function retrieves a list of ingredients and is designed to be
 * the canonical API for the underlying data storage layer.
 */
func (f *Fetcher) AllIngredients() []proto.Ingredient {
	ingredients := make([]proto.Ingredient, 0)
	iter := f.SS.Database.Ingredients.Find(nil).Iter()

	result := proto.Ingredient{}
	for iter.Next(&result) {
		ingredients = append(ingredients, result)
	}

	return ingredients
}
