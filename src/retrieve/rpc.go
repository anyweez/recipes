package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	gproto "code.google.com/p/goprotobuf/proto"
	"log"
	"io/ioutil"
	"net/http"
	proto "proto"
	"strings"
)

type Retriever int

type IngredientList struct {
	Ingredients		[]string
}

type GraphResult struct {
	Result	[]GraphNode
}

type GraphNode struct {
	Id		string
}

/**
 * GetPartialRecipes fetches a list of recipes that contain all of the
 * ingredients provided in the input IngredientList.
 */
func (r *Retriever) GetPartialRecipes(il *IngredientList, reply *proto.RecipeBook) error {
	log.Println("RPC REQUEST:" + strings.Join(il.Ingredients, ","))
	url := fmt.Sprintf("http://%s/api/v1/query/gremlin", *OUTPUT_QUADS)

	recipes := make(map[string]bool, 0)
	
	for _, ingredient := range il.Ingredients {
		// Body (Gremlin query)
		data := []byte(fmt.Sprintf("g.Vertex(\"%s\").In(\"contains\").All()", ingredient))
		resp, _ := http.Post( url, "text/plain", bytes.NewReader(data) )

		defer resp.Body.Close()
	
		bd, _ := ioutil.ReadAll(resp.Body)
		gr := GraphResult{}
		json.Unmarshal(bd, &gr)
		
		// Create a set.
		for _, node := range gr.Result {
			recipes[node.Id] = true
		}
	}
	
	for rid := range recipes {
		recipe := proto.Recipe{
			Id: gproto.String(rid),
		}
		
		reply.Recipes = append(reply.Recipes, &recipe)
	}
		
	return nil
}


