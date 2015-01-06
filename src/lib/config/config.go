package config

import (
	"fmt"
)

type RecipesConfig struct {
	MongoAddress 	string
	MongoPort    	int
	FreebaseDump	string
}

func New(filename string) RecipesConfig {
	return RecipesConfig{
		MongoAddress: 	"historian",
		MongoPort:    	27017,
		FreebaseDump:	"/mnt/vortex/corpora/freebase/freebase.all",
	}
}

func (c *RecipesConfig) Mongo() string {
	return fmt.Sprintf("%s:%d", c.MongoAddress, c.MongoPort)
}
