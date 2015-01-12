package recipes

import (
	"labix.org/v2/mgo/bson"
)

type PageRecord struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page []byte
	// The HTML content of the page.
	Content []byte
}

/*
func GetAll() []proto.Recipes {

}

func GetAllRaw(conf config.RecipesConfig) [][]byte {

}
*/
