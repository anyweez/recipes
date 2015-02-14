package fetch

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	proto "proto"
)

type PageRecord struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page []byte
	// The HTML content of the page.
	Content []byte
}

var rc *mgo.Collection

func init() {
	conf, _ := config.New("recipes.conf")

	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("[fetch/recipes] Recipe retrieval API can't connect to MongoDB instance: " + conf.Mongo.ConnectionString())
	}

	rc = session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
}

func Recipe(recipeId string) proto.Recipe {
	recipe := proto.Recipe{}
	rc.Find(bson.M{"id": recipeId}).One(&recipe)

	return recipe
}

/**
 * Fetch all recipes. Note that this function retrieves all in batch and
 * does NOT return an iterator. If there are lots of recipes (too many
 * to comfortably store in RAM) then this calling this function may be
 * a bad idea.
 */
func AllRecipes() []proto.Recipe {
	recipes := make([]proto.Recipe, 0)
	iter := rc.Find(nil).Iter()

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
