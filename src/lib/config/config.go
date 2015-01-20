package config

import (
	"fmt"
	gcfg "code.google.com/p/gcfg"
	"log"
)

type RecipesConfig struct {
	Mongo		MongoConfig
	Freebase	FreebaseConfig
}

type MongoConfig struct {
	Address						string
	Port						int
	DatabaseName				string
	RawCollection				string
	RecipeCollection			string
	IngredientCollection		string
	RecommendationCollection	string
}

type FreebaseConfig struct {
	DumpLocation	string
}

func New(filename string) RecipesConfig {
	c := RecipesConfig{}
	// Read the configuration.
	err := gcfg.ReadFileInto(&c, filename)
	if err != nil {
		log.Fatal("Error reading configuration file: " + err.Error())
	}

	return c
}

func (mc *MongoConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%d", mc.Address, mc.Port)
}
