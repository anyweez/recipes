package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
//	gproto "code.google.com/p/goprotobuf/proto"
	"io/ioutil"
	"lib/fetch"
	"log"
	"math/rand"
	"net/http"
	proto "proto"
	"strings"
)

type Retriever int

type IngredientList struct {
	Ingredients []string
}

type GraphResult struct {
	Result []GraphNode
}

type GraphNode struct {
	Id string
}

/**
 * The request object that contains all of the fields required to make a
 * best recipe RPC call.
 * 
 * All fields are currently required.
 */
type BestRecipesRequest struct {
	// Will eventually be removed from the API. Currently useful for testing
	Seed		int64
	UserId		uint64
	GroupId		uint64
	
	// The number of recipes desired, if possible (not guaranteed)
	Count		int
}

/**
 * GetPartialRecipes fetches a list of recipes that contain all of the
 * ingredients provided in the input IngredientList.
 */
func (r *Retriever) GetPartialRecipes(il *IngredientList, reply *proto.RecipeBook) error {
	//session, _ := mgo.Dial(*MONGO)
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB instance.")
	}
	c := session.DB("recipes").C("parsed")

	log.Println("RPC REQUEST:" + strings.Join(il.Ingredients, ","))
	url := fmt.Sprintf("http://%s/api/v1/query/gremlin", *OUTPUT_QUADS)

	recipes := make(map[string]int, 0)

	for _, ingredient := range il.Ingredients {
		// Body (Gremlin query)
		data := []byte(fmt.Sprintf("g.Vertex(\"%s\").In(\"contains\").All()", ingredient))
		resp, err := http.Post(url, "text/plain", bytes.NewReader(data))

		if err != nil {
			log.Fatal("Couldn't update Cayley: " + err.Error())
		}

		defer resp.Body.Close()

		bd, _ := ioutil.ReadAll(resp.Body)
		gr := GraphResult{}
		json.Unmarshal(bd, &gr)

		// Create a set.
		for _, node := range gr.Result {
			_, exists := recipes[node.Id]
			if exists {
				recipes[node.Id] += 1
			} else {
				recipes[node.Id] = 1
			}
		}
	}

	for key, val := range recipes {
		log.Println(fmt.Sprintf("%s => %d", key, val))
		// Keep the recipe if it was retrieved for each ingredient.
		if val == len(il.Ingredients) {
			recipe := proto.Recipe{}
			c.Find(bson.M{"id": key}).One(&recipe)

			reply.Recipes = append(reply.Recipes, &recipe)
		}
	}

	return nil
}

func (r *Retriever) GetIngredients(na string, ingr *[]proto.Ingredient) error {
	for _, in := range fetch.AllIngredients() {
		*ingr = append(*ingr, in)
	}

	return nil
}

/**
 * Get the top recommended recipes for a specified user in a specified group.
 * This function will pull first from the recommended recipe cache (if
 * available) and fill the remaining open slots with randomly selected
 * recipes (lightly filtered).
 */
func (r *Retriever) GetBestRecipes(request BestRecipesRequest, recipes *[]proto.Recipe) error {
	log.Println( fmt.Sprintf("Request count=%d for user=%d, group=%d", request.Count, request.UserId, request.GroupId) )
	// Seed the random number generator
	rand.Seed(request.Seed)

	// Connect to Mongo
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	defer session.Close()
	
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB instance.")
	}
	
	// Get as many recommended recipes as possible (up to request.Count).
	rset := fetchRecommended(session, request, request.Count)
	
	// If we didn't find as many recipes as requested, find some more (probably
	// lower quality recipes) to backfill with and merge the two lists together.
	//
	// Note that this is likely slower than fetchRecommended so the goal is
	// to make it relatively unlikely that this branch is required. The main
	// lever we have for that is the size of the recommended recipes cache.
	if len(rset) < request.Count {
		rset = merge(rset, fetchMore(session, request, request.Count - len(rset)))
	}	

	for _, recipe := range rset {
		// TODO: copy over to recommended if it doesn't exist and update ServingStatus 
		*recipes = append(*recipes, recipe)
	}
	return nil
}

/**
 * Fetch recommended recipes from a precomputed datastore. The list is simply treated
 * like a queue and recipes are read from the front of the list.
 * 
 * Once read, recipes are (potentially?) labeled as "returned" until they're answered,
 * at which point they're labeled as "answered."
 */
func fetchRecommended(session *mgo.Session, request BestRecipesRequest, count int) []proto.Recipe {
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecommendationCollection)

	// TODO: filter out recipes where Recipe.ServingRecord.User.Id = userid on the database.
	iter := c.Find(bson.M{"group_id": request.GroupId}).Iter()
	recipe := proto.Recipe{}
	recipes := make([]proto.Recipe, 0)
	
	// Iterate through all returned records until either we run out of records or
	// we get all that we came for.
	for iter.Next(&recipe) || len(recipes) == count {
		eligible := true
		
		// Check to see if this record has been shown to a user before. If not, include it
		// in the returned set.
		for _, sr := range recipe.ServingRecord {
			if *sr.User.Id == request.UserId && *sr.Status != proto.Recipe_ServingRecord_NOT_RETURNED {
				eligible = false
			}
		}
		
		if eligible {
			recipes = append(recipes, recipe)
		}
	}

	return recipes
}

/**
 * Fetch random recipes to backfill for a shortage from fetchRecommended()
 * and do some light quality filtering.
 */
func fetchMore(session *mgo.Session, request BestRecipesRequest, count int) []proto.Recipe {
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
	fmt.Println( fmt.Sprintf("Fetching %d recipes.", count) )
	// TODO: filter out recipes where Recipe.ServingRecord.User.Id = userid on the database.
	// TODO: is `context` the right word in Mongo-speak?
	query := c.Find(nil)
	numRecords, _ := query.Count()
	
	recipe := proto.Recipe{}
	recipes := make([]proto.Recipe, 0)

	// Randomly select COUNT values, which will be indeces of recipes
	// that we want to use.
	chosen := make([]int, 0, count+1)

	for i := 0; i < count+1; i++ {
		index := rand.Int() % numRecords
		chosen = append(chosen, index)
	}
	
	// Skip to the first record and grab it.
	i := 0
	query.Skip(chosen[i])
	err := query.One(&recipe)

	for err == nil && len(recipes) < count {
		// Store this recipe.
		// TODO: do some basic and quick quality checks
		fmt.Println(*recipe.Name)
		recipes = append(recipes, recipe)

		// Retrieve a new record
		i++
		query.Skip(chosen[i])
		err = query.One(&recipe)	
	}

	return recipes
}

func merge(first []proto.Recipe, second []proto.Recipe) []proto.Recipe {
	// Make a new slice that is at least as long as the first list.
	recipes := make([]proto.Recipe, 0, len(first))
	
	// Copy the first list in, then copy over anything from the second list
	// that's unique from elements in the first list.
	for _, r := range first {
		recipes = append(recipes, r)
	}
	
	for _, r := range second {
		unique := true
		
		for _, existing := range recipes {
			// If this recipe exists in the list, don't keep it anymore.
			if r.Id == existing.Id {
				unique = false
				break
			}
		}
		
		// If the recipe doesn't already exist in the list, add it.
		if unique {
			recipes = append(recipes, r)
		}
	}
	
	return recipes
}
