package config

import (
	"fmt"
	gcfg "code.google.com/p/gcfg"
)

type RecipesConfig struct {
	Mongo		MongoConfig
	Freebase	FreebaseConfig
}

type MongoConfig struct {
	Address					string
	Port					string
	DatabaseName			string
	RawCollection			string
	RecipeCollection		string
	IngredientCollection	string
}

type FreebaseConfig struct {
	DumpLocation	string
}

func New(filename string) RecipesConfig {
	c := RecipesConfig{}
	// Read the configuration.
	gcfg.ReadFileInto(&c, filename)

	return c
}

func (c *RecipesConfig) Mongo() string {
	return fmt.Sprintf("%s:%d", c.MongoAddress, c.MongoPort)
}
