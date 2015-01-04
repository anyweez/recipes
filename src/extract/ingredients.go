package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"lib/config"	
	"log"
	gproto "code.google.com/p/goprotobuf/proto"
	proto "proto"
)

/**
 * This function parses a Freebase archive and generates a list of
 * Ingredient structures.
 */
func ExtractIngredients(conf config.RecipesConfig) []proto.Ingredient {
	ingredients := make([]proto.Ingredient, 0)
	
	ingredients = append(ingredients, proto.Ingredient{ Name: gproto.String("tater salad") })
	ingredients = append(ingredients, proto.Ingredient{ Name: gproto.String("cole slaw") })
	
	return ingredients
}

/**
 * This function updates a MongoDB instance with ingredient data (usually
 * extracted from ExtractIngredients but not necessarily). It creates new
 * records if they don't exist, updates them if they do, and leaves them
 * untouched if they're not present in the input slice.
 */
func UpdateIngredients(conf config.RecipesConfig, ingredients []proto.Ingredient) error {
	session, err := mgo.Dial( fmt.Sprintf("%s:%d", conf.MongoAddress, conf.MongoPort) )
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB to update ingredient list: " + err.Error())
	}
	c := session.DB("recipes").C("ingredients")

	for _, ingr := range ingredients {
		c.Insert(ingr)
	}
	
	return nil
}
