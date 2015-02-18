package fetch

import (
	"labix.org/v2/mgo/bson"
	proto "proto"
)

type PageRecord struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page []byte
	// The HTML content of the page.
	Content []byte
}

func (f *Fetcher) Recipe(recipeId string) proto.Recipe {
	recipe := proto.Recipe{}
	f.SS.Database.Recipes.Find(bson.M{"id": recipeId}).One(&recipe)

	return recipe
}

/**
 * Fetch all recipes. Note that this function retrieves all in batch and
 * does NOT return an iterator. If there are lots of recipes (too many
 * to comfortably store in RAM) then this calling this function may be
 * a bad idea.
 */
func (f *Fetcher) AllRecipes() []proto.Recipe {
	recipes := make([]proto.Recipe, 0)
	iter := f.SS.Database.Recipes.Find(nil).Iter()

	recipe := proto.Recipe{}
	for iter.Next(&recipe) {
		recipes = append(recipes, recipe)
	}

	return recipes
}

/*
func GetAllRaw(conf config.RecipesConfig) [][]byte {

}
*/
