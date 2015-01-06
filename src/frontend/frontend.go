package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	retrieve "retrieve"
	proto "proto"
	"strings"
)

var INGREDIENTS = flag.String("ingredients", "m/0ggm5yy", "The ingredients we should search for.")
var RETRIEVER = flag.String("retriever", "127.0.0.1:14501", "")
var MONGO = flag.String("mongo", "127.0.0.1:27017", "")

func list_ingredients(w http.ResponseWriter, r *http.Request) {
	// Create a connection to the RPC server to handle this request.
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}

	ingredients := make([]proto.Ingredient, 0)
	err = client.Call("Retriever.GetIngredients", "hi", &ingredients)
	if err != nil {
		log.Fatal("Cannot retrieve ingredient list.")
	}

	data, _ := json.Marshal(ingredients)
	fmt.Fprintf(w, string(data))
}

func find_recipes(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}

   	var il retrieve.IngredientList
   	il.Ingredients = make([]string, 0)
   	il.Ingredients = append(il.Ingredients, strings.Split("m/0ggm5yy", ",")...)

   	rb := proto.RecipeBook{}

   	err = client.Call("Retriever.GetPartialRecipes", il, &rb)
	data, _ := json.Marshal(rb)
	
	fmt.Fprintf(w, string(data))
}

func main() {
	flag.Parse()

	http.HandleFunc("/api/ingredients", list_ingredients)
	http.HandleFunc("/api/recipes", find_recipes)
	// No-op handler for favicon.ico, since it'll otherwise generate an extra call to index.
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	// Serve any files in static/ directly from the filesystem.
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL.Path[1:])
		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
	})

	log.Println("Awaiting requests...")
	log.Fatal("Couldn't listen on port 8088:", http.ListenAndServe(":8088", nil))

	/*

	   	var il retrieve.IngredientList
	   	il.Ingredients = make([]string, 0)

	   	il.Ingredients = append(il.Ingredients, strings.Split(*INGREDIENTS, ",")...)
	   	rb := proto.RecipeBook{}
	   //	rb.Recipes = make([]proto.Recipe, 0)

	   	err = client.Call("Retriever.GetPartialRecipes", il, &rb)
	   	for i, recipe := range rb.Recipes {
	   		fmt.Println( fmt.Sprintf("  %d. %s (%s)", i+1, *recipe.Name, *recipe.Id) )
	   	}
	   	* */
}
