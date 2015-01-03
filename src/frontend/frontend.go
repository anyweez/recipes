package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	retrieve "retrieve"
	proto "proto"
	"strings"
)

var INGREDIENTS = flag.String("ingredients", "m/0ggm5yy", "The ingredients we should search for.")
var RETRIEVER = flag.String("retriever", "127.0.0.1:14501", "")
var MONGO = flag.String("mongo", "127.0.0.1:27017", "")

func main() {
	flag.Parse()
	
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}
	
	var il retrieve.IngredientList
	il.Ingredients = make([]string, 0)
	
	il.Ingredients = append(il.Ingredients, strings.Split(*INGREDIENTS, ",")...)
	
	rb := proto.RecipeBook{}
//	rb.Recipes = make([]proto.Recipe, 0)
	
	err = client.Call("Retriever.GetPartialRecipes", il, &rb) 
	for i, recipe := range rb.Recipes {
		fmt.Println( fmt.Sprintf("  %d. %s (%s)", i+1, *recipe.Name, *recipe.Id) )
	}
}
