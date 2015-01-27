package config

import (
	"errors"
	"fmt"
	gcfg "code.google.com/p/gcfg"
)

type RecipesConfig struct {
	Mongo		MongoConfig
	Freebase	FreebaseConfig
	Rpc			RPCConfig
}

type MongoConfig struct {
	Address						string
	Port						int
	DatabaseName				string
	RawCollection				string
	RecipeCollection			string
	IngredientCollection		string
	RecommendationCollection	string
	ResponseCollection			string
	UserCollection				string
	GroupCollection				string
}

type RPCConfig struct {
	Address						string
	Port						int
}

type FreebaseConfig struct {
	DumpLocation	string
}

func New(filename string) (RecipesConfig, error) {
	c := RecipesConfig{}
	// Read the configuration.
	err := gcfg.ReadFileInto(&c, filename)
	if err != nil {
		return c, errors.New("Error reading configuration file: " + err.Error())
	}

	return c, nil
}

func (mc *MongoConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%d", mc.Address, mc.Port)
}

func (rpcc *RPCConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%d", rpcc.Address, rpcc.Port)	
}
