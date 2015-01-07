package main

import (
	"bufio"
	"fmt"
	"labix.org/v2/mgo"
	"lib/config"	
	"log"
	gproto "code.google.com/p/goprotobuf/proto"
	"os"
	proto "proto"
	"strings"
)

/**
 * This function takes a line of a RDF file and breaks it up into pieces.
 * It also strips quotation marks, removes angled brackets, etc.
 */
func split(line string) (string, string, string) {
	parts := strings.SplitN(line, " ", 3)

	// Cut the last two characters off of the final part (" .")
	return convertFreebaseId(parts[0]), convertFreebaseId(parts[1]), convertFreebaseId(parts[2][0:len(parts[2])-2])
}

func convertFreebaseId(uri string) string{
     if strings.HasPrefix(uri, "<") && strings.HasSuffix(uri, ">") {
        var id = uri[1 : len(uri)-1]
        id = strings.Replace(id, "http://rdf.freebase.com/ns", "", -1)
        id = strings.Replace(id, ".", "/", -1)
        return id
     }
     
     return uri
 }

func store(mapping map[string]*proto.Ingredient, subj string, pred string, obj string) bool {
	// This switch statement identifies which predicates should be stored,
	// and what field they should be stored in.
	switch (pred) {
		case "/type/object/name":
			ingredient, _ := mapping[subj]
			ingredient.Name = gproto.String(obj)
			
			return true
			break
		default:
			break
	}
	
	return false
}

/**
 * Checks whether a given tuple invalidates the subject. Note that the
 * subject only needs to be invalidated (marked non-keeper) once in order
 * to be skipped.
 */
func isKeeper(subj, pred, obj string) bool {
	// If notable_type is /food/food (expecting 8,615 cases at time of writing
	// according to http://www.freebase.com/food?schema=).
	if (pred == "/common/topic/notable_types" && obj == "/m/05yxcqj") {
		return true
	}

	return false
}

/**
 * This function parses a Freebase archive and generates a list of
 * Ingredient structures.
 */
func ExtractIngredients(conf config.RecipesConfig) []*proto.Ingredient {
	ingredients := make([]*proto.Ingredient, 0)
	
	// Step 1: open file w/ reader (note that it can be VERY big so it needs to be buffered)
	fp, err := os.Open(conf.Freebase.DumpLocation)
	if err != nil {
		log.Fatal("Couldn't open Freebase sample file at " + conf.Freebase.DumpLocation)
	}
	
	scanner := bufio.NewScanner( bufio.NewReader(fp) )
	
	// Step 2: create two maps, one that keeps track of ingredient data (keyed by mid)
	//   and the other that keeps track of whether a given ingredient is a "keeper."
	//   Note that all ingredients are assuming to be keepers until a property is
	//   discovered that makes them no longer interesting (such as being of the wrong
	//   type).
	im := make(map[string]*proto.Ingredient)
	iv := make(map[string]bool)

	current := &proto.Ingredient{}
	current_mid := ""
	line_count := 0
	log.Println("Starting scan...")
	for scanner.Scan() {
		subj, pred, obj := split(scanner.Text())
		
		// If this is a new mid, either get rid of or store current.
		if subj != current_mid {
			_, keeper := iv[current_mid]

			// If the mid exists and is a keeper, store it. Make sure we know
			// the name or it's not interesting.
			if keeper && im[current_mid].Name != nil {
				log.Println( fmt.Sprintf("Keeping %s!", *im[current_mid].Name) )
				ingredients = append(ingredients, im[current_mid])
			// Or forget that this key ever existed.
			} else {
				delete(iv, current_mid)
				delete(im, current_mid)
			}
			
			im[subj] = &proto.Ingredient{}
			current = im[subj]
			current.Ingrids = append( current.Ingrids, subj )
			current_mid = subj
		}
		
		// Check whether this record indicates that this is a record to keep.
		if isKeeper(subj, pred, obj) {
			iv[subj] = true
			log.Println( fmt.Sprintf("%s is a keeper!", subj) )
		}
		
		// Store the field on the current object
		store(im, subj, pred, obj)		
		line_count += 1
	}

	log.Println( fmt.Sprintf("%d lines read.", line_count) )
	return ingredients
}

/**
 * This function updates a MongoDB instance with ingredient data (usually
 * extracted from ExtractIngredients but not necessarily). It creates new
 * records if they don't exist, updates them if they do, and leaves them
 * untouched if they're not present in the input slice.
 */
func UpdateIngredients(conf config.RecipesConfig, ingredients []proto.Ingredient) error {
	session, err := mgo.Dial( conf.Mongo.ConnectionString() )
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB to update ingredient list: " + err.Error())
	}
	c := session.DB("recipes").C("ingredients")

	for _, ingr := range ingredients {
		c.Insert(ingr)
	}
	
	return nil
}
