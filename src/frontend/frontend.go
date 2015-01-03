package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	retrieve "retrieve"
	proto "proto"
)

var INGREDIENTS = flag.String("ingredients", "m/0ggm5yy", "The ingredients we should search for.")
var RETRIEVER = flag.String("retriever", "127.0.0.1:14501", "")

func main() {
	flag.Parse()
	
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}
	
	var il retrieve.IngredientList
	il.Ingredients = make([]string, 0)
	il.Ingredients = append(il.Ingredients, "m/0ggm5yy")
	
	rb := proto.RecipeBook{}
//	rb.Recipes = make([]proto.Recipe, 0)
	
	err = client.Call("Retriever.GetPartialRecipes", il, &rb) 
	for i, recipe := range rb.Recipes {
		fmt.Println( fmt.Sprintf("  %d. %s", i+1, *recipe.Id) )
	}
}
