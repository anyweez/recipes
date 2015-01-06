package config

import (
	"fmt"
)

type RecipesConfig struct {
	MongoAddress 				string
	MongoPort    				int
	FreebaseDump				string
	MongoDatabase				string
	MongoRawCollection			string
	MongoRecipeCollection		string
	MongoIngredientCollection	string	
}

func New(filename string) RecipesConfig {
	return RecipesConfig{
		MongoAddress: 				"historian",
		MongoPort:    				27017,
		MongoDatabase:				"recipes",
		MongoRawCollection:			"scraped",
		MongoRecipeCollection:		"recipes",
		MongoIngredientCollection:	"ingredients",	
		FreebaseDump:				"/mnt/vortex/corpora/freebase/freebase.all",
	}
}

func (c *RecipesConfig) Mongo() string {
	return fmt.Sprintf("%s:%d", c.MongoAddress, c.MongoPort)
}
