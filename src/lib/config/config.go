package config

import (
	gcfg "code.google.com/p/gcfg"
	"errors"
	"fmt"
	"os"
)

type RecipesConfig struct {
	Mongo    MongoConfig
	Freebase FreebaseConfig
	Rpc      RPCConfig
	Frontend FrontendConfig
}

type MongoConfig struct {
	Address                  string
	Port                     int
	DatabaseName             string
	RawCollection            string
	RecipeCollection         string
	IngredientCollection     string
	RecommendationCollection string
	ResponseCollection       string
	UserCollection           string
	GroupCollection          string
	MealsCollection          string
	VotesCollection          string
}

type RPCConfig struct {
	Address string
	Port    int
}

type FreebaseConfig struct {
	DumpLocation string
}

type FrontendConfig struct {
	Port int
}

func New(filename string) (RecipesConfig, error) {
	c := RecipesConfig{}
	// Read the configuration.
	path, _ := os.Getwd()
	fmt.Println(path)

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
