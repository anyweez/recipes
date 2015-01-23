package main

/**
 * This package contains the core retrieval functions for the application,
 * which are exposed through the RPC interface defined in retrieve.go in
 * this package.
 * 
 * These RPC's are specifically online retrieval-based functions and depend
 * on offline processes to do a bunch of quality work.
 */

import (
	"bytes"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	gproto "code.google.com/p/goprotobuf/proto"
	"io/ioutil"
	"lib/config"
	"lib/fetch"
	log "logging"
	"math/rand"
	"net/http"
	proto "proto"
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
 * EXPOSED AS RPC
 * 
 * Returns a list of all known ingredients, to be used for populating lists,
 * autocomplete, etc.
 */
func (r *Retriever) GetIngredients(na string, ingr *[]proto.Ingredient) error {
	for _, in := range fetch.AllIngredients() {
		*ingr = append(*ingr, in)
	}

	return nil
}

/**
 * EXPOSED AS RPC
 * 
 * Get the top recommended recipes for a specified user in a specified group.
 * This function will pull first from the recommended recipe cache (if
 * available) and fill the remaining open slots with randomly selected
 * recipes (lightly filtered).
 */
func (r *Retriever) GetBestRecipes(request BestRecipesRequest, recipes *[]proto.Recipe) error {
	// Create a stable request ID for this RPC call.
	le := log.New("GetBestRecipe", log.Fields{
		"userid": request.UserId,
		"groupid": request.GroupId,
	})
	
	// Seed the random number generator
	rand.Seed(request.Seed)

	// Connect to Mongo
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	defer session.Close()
	
	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't connect to MongoDB instance.", log.Fields{
			"db": conf.Mongo.DatabaseName,
			"ip": conf.Mongo.Address,
			"port": conf.Mongo.Port,
		})
	}
	
	// Get as many recommended recipes as possible (up to request.Count).
	rset := fetchRecommended(session, request, request.Count)
	numRecommended := len(rset)
	
	// If we didn't find as many recipes as requested, find some more (probably
	// lower quality recipes) to backfill with and merge the two lists together.
	//
	// Note that this is likely slower than fetchRecommended so the goal is
	// to make it relatively unlikely that this branch is required. The main
	// lever we have for that is the size of the recommended recipes cache.
	if numRecommended < request.Count {
		rset = merge(rset, fetchMore(session, request, request.Count - len(rset)))
	}	

	for _, recipe := range rset {
		// TODO: copy over to recommended if it doesn't exist and update ServingStatus 
		*recipes = append(*recipes, recipe)
	}

	// Log the completion of the event.
	le.Update(log.STATUS_COMPLETE, "", log.Fields{
		"recommended": numRecommended,
		"more": len(rset) - numRecommended, 		
	})
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
 * and do some light quality filtering. Recipes returned from this function
 * can be returned multiple times (a ServingRecord is not maintained for them).
 */
func fetchMore(session *mgo.Session, request BestRecipesRequest, count int) []proto.Recipe {
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
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

/**
 * EXPOSED AS RPC
 * 
 * Fetch the list of all recent responses for a given recipe within a group.
 */
func (r *Retriever) GetRecipeResponse(request RecipeResponseRequest, response *[]proto.RecipeResponses_RecipeResponse) error {
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.ResponseCollection)
	responses := proto.RecipeResponses{}
	c.Find(bson.M{"group_id": request.GroupId}).One(&responses)
	
	// If the recipe has any responses, find them and add them.
	for _, rr := range responses.Responses {
		if *rr.Recipe.Id == request.RecipeId {
			*response = append(*response, *rr)
		}
	}
	
	return nil
}

/**
 * EXPOSED AS RPC 
 * 
 * Record a user's response to a recipe. Note that all responses are made in
 * the context of a group, so a single call to this function will only
 * store the answer once; if it should be stored for all of the user's
 * groups then the call will need to be made multiple times.
 * 
 * TODO: this function currently contains a race condition if members of
 * the same group submit responses in close proximity. Need to use atomic update
 * and commit.
 */
func (r *Retriever) PostRecipeResponse(request RecipeResponse, success *bool) error {
	// Create a stable logging request for this RPC call.
	le := log.New("PostRecipeResponse", log.Fields{
		"userid": request.UserId,
		"groupid": request.GroupId,
		"recipeid": request.RecipeId,		
	})

	// Connect to Mongo
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.ResponseCollection)
	// Atomically fetch proto.RecipeResponses object by request.GroupId,
	// and add this respones as another Response.
	user := fetch.User(request.UserId)
	recipe := fetch.Recipe(request.RecipeId)
	resp_enum := proto.RecipeResponses_RecipeResponse_NO
	if request.Response {
		resp_enum = proto.RecipeResponses_RecipeResponse_YES
	}
	
	cnt, err := c.Find(bson.M{"group_id": request.GroupId}).Count()
	if err != nil {
		le.Update(log.STATUS_ERROR, err.Error(), nil)
	}
	
	resp := proto.RecipeResponses_RecipeResponse {
		User: &user,
		Response: &resp_enum,
		Recipe: &recipe,
	}
	
	// Append to the existing record.
	// TODO: simplify this by using an upsert?
	if cnt > 0 {
		rr := proto.RecipeResponses{}
		c.Find(bson.M{"group_id": request.GroupId}).One(&rr)
			
		rr.Responses = append(rr.Responses, &resp)
		c.Update(bson.M{"group_id": request.GroupId}, rr)
	// Create a new record.
	} else {
		rr := proto.RecipeResponses{
			GroupId: gproto.Uint64(request.GroupId),
		}
		
		rr.Responses = append(rr.Responses, &resp)
		c.Insert(rr);
	}
		
	le.Update(log.STATUS_COMPLETE, "", nil)
	
	return nil	
}

/**
 * EXPOSED AS RPC
 * 
 * GetPartialRecipes fetches a list of recipes that contain all of the
 * ingredients provided in the input IngredientList.
 */
func (r *Retriever) GetPartialRecipes(il *IngredientList, reply *proto.RecipeBook) error {
	log.Info("Inbound RPC request", log.Fields{
		"rpc": "GetPartialRecipes",
	})
	
	conf, _ := config.New("recipes.conf")
	//session, _ := mgo.Dial(*MONGO)
	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB instance.", log.Fields{
			"db": conf.Mongo.DatabaseName,
			"ip": conf.Mongo.Address,
			"port": conf.Mongo.Port,
		})	
	}
	c := session.DB("recipes").C("parsed")

	url := fmt.Sprintf("http://%s/api/v1/query/gremlin", *OUTPUT_QUADS)
	recipes := make(map[string]int, 0)

	for _, ingredient := range il.Ingredients {
		// Body (Gremlin query)
		data := []byte(fmt.Sprintf("g.Vertex(\"%s\").In(\"contains\").All()", ingredient))
		resp, err := http.Post(url, "text/plain", bytes.NewReader(data))

		if err != nil {
			log.Fatal("Couldn't update Cayley: " + err.Error(), nil)	
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
		// Keep the recipe if it was retrieved for each ingredient.
		if val == len(il.Ingredients) {
			recipe := proto.Recipe{}
			c.Find(bson.M{"id": key}).One(&recipe)

			reply.Recipes = append(reply.Recipes, &recipe)
		}
	}

	return nil
}
