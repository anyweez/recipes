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
 * This function takes a line of a RDF file and breaks it up into pieces.
 * It also strips quotation marks, removes angled brackets, etc.
 */
func split(line string) (string, string, string) {
	return "one", "two", "three"
}

func store(mapping *map[string]bool, subj string, pred string, obj string) {
	// This switch statement identifies which predicates should be stored,
	// and what field they should be stored in.
	switch (pred) {
		case "/name":
			mapping[subj].Name = obj
			break
		default:
			break
	}
}

/**
 * Checks whether a given tuple invalidates the subject. Note that the
 * subject only needs to be invalidated (marked non-keeper) once in order
 * to be skipped.
 */
func isKeeper(subj, pred, obj string) bool {
	return false
}

/**
 * This function parses a Freebase archive and generates a list of
 * Ingredient structures.
 */
func ExtractIngredients(conf config.RecipesConfig) []proto.Ingredient {
	ingredients := make([]proto.Ingredient, 0)
	
	// Step 1: open file w/ reader (note that it can be VERY big so it needs to be buffered)
	
	// Step 2: create two maps, one that keeps track of ingredient data (keyed by mid)
	//   and the other that keeps track of whether a given ingredient is a "keeper."
	//   Note that all ingredients are assuming to be keepers until a property is
	//   discovered that makes them no longer interesting (such as being of the wrong
	//   type).
	im := make(map[string]proto.Ingredient)
	iv := make(map[string]bool)
	
##	for _, line := line of freebase file {
		subj, pred, obj := split(line)
		
		// Check whether this record invalidates the subject.
		if !isKeeper(subj, pred, obj) {
			iv[subj] = false
##			delete im[subj]
		} else {
			// If not, let's see if its a fact that needs to be saved.
			valid, exists := iv[subj]
			// Save the record if so, otherwise just skip it.
			if exists && valid {
				store(&im, subj, pred, obj)
			}
		}			
	}
	
	// Convert the map into a slice that can be returned.
	for _, ingr := range im {
		ingredients = append(ingredients, ingr)
	}
	
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
