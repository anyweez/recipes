package main

import (
	"encoding/json"
	"fmt"
	log "logging"
	"math/rand"
	"net/http"
	"net/rpc"
	proto "proto"
	retrieve "retrieve"
	"strconv"
	"time"
)

func best_recipes(w http.ResponseWriter, r *http.Request) {
	le := log.New("web_request", log.Fields{
		"handler": "best_recipes",
	})

	client, err := rpc.DialHTTP("tcp", conf.Rpc.ConnectionString())
	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't connect to retriever: "+err.Error(), nil)
	}

	rand.Seed(time.Now().UnixNano())
	seed := int64(rand.Int())

	rseed_str := r.URL.Query().Get("seed")
	rseed, err := strconv.ParseInt(rseed_str, 10, 64)

	if err == nil {
		seed = rseed
	}

	recipes := make([]proto.Recipe, 0)
	err = client.Call("Retriever.GetBestRecipes", retrieve.BestRecipesRequest{
		Seed:    seed,
		UserId:  5,
		GroupId: 1,
		Count:   5,
	}, &recipes)

	data, _ := json.Marshal(recipes)
	fmt.Fprintf(w, string(data))

	le.Update(log.STATUS_COMPLETE, "", nil)
}
