package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"math/rand"
	"encoding/json"
	proto "proto"
	"strconv"
	"time"
)

func best_recipes(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}
	
	rand.Seed(time.Now().Unix())
	seed := int64(rand.Int())

	rseed_str := r.URL.Query().Get("seed")
	rseed, err := strconv.ParseInt(rseed_str, 10, 64)
	
	if err == nil {
		seed = rseed
	}

	log.Println( fmt.Sprintf("Seed: %d", seed) )

	recipes := make([]proto.Recipe, 0)
   	err = client.Call("Retriever.GetBestRecipes", seed, &recipes)
   	
	data, _ := json.Marshal(recipes)
	fmt.Fprintf(w, string(data))
}
