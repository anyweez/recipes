package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	gproto "code.google.com/p/goprotobuf/proto"
	"io/ioutil"
	"lib/ingredients"
	"log"
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
	for _, in := range ingredients.GetAll() {
		*ingr = append(*ingr, in)
	}

	return nil
}
