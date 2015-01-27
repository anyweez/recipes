package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	proto "proto"
)

func list_ingredients(w http.ResponseWriter, r *http.Request) {
	// Create a connection to the RPC server to handle this request.
	client, err := rpc.DialHTTP("tcp", conf.Rpc.ConnectionString())
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
