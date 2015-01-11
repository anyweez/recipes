package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"encoding/json"
	proto "proto"
)

func best_recipes(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}

	recipes := make([]proto.Recipe, 0)
   	err = client.Call("Retriever.GetBestRecipes", 12, &recipes)
   	
	data, _ := json.Marshal(recipes)

	fmt.Fprintf(w, string(data))
}
