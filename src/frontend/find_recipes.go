package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	retrieve "retrieve"
	"encoding/json"
	proto "proto"
	"strings"
)

func find_recipes(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", conf.Rpc.ConnectionString())
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
