package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"lib/config"
	log "logging"
	"net/http"
	"net/rpc"
	retrieve "retrieve"
	router "frontend/router"
	"strconv"
)

var conf config.RecipesConfig

func init() {
	// TODO: check that this was read in correctly.
	conf, _ = config.New("recipes.conf")
	
//	if err != nil {
//		le := log.New("init", nil)
//		le.Update(log.STATUS_FATAL, err.Error(), nil)
//	}
}

// Fetch the index page for the "rate" URL path.
func rate_index_handler(w http.ResponseWriter, r *http.Request) {
	le := log.New("web_request", log.Fields{
		"handler": "rate_index_handler",
	})

	data, err := ioutil.ReadFile("html/rate/index.html")
	if err != nil {
		le.Update(log.STATUS_FATAL, "rate/index.html not present!", nil)
		http.NotFound(w, r)
	} else {
		w.Write(data)
	}
	
	le.Update(log.STATUS_COMPLETE, "", nil)
}

func submit_response(w http.ResponseWriter, r *http.Request) {
	le := log.New("web_request", log.Fields{
		"handler": "submit_response",
	})
	
	client, err := rpc.DialHTTP("tcp", conf.Mongo.ConnectionString())
	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't connect to retriever: " + err.Error(), nil)
	}

	recipe_id := r.URL.Query().Get("recipe")
	response, err := strconv.ParseBool(r.URL.Query().Get("response")) // "true" or "false"
	if err != nil {
		le.Update(log.STATUS_ERROR, "Invalid value for `response`; must be 'true' or 'false'.", nil)
		// TODO: return something to the user, or at least handle this more elegantly.
		return
	}
	
	success := true
	
   	err = client.Call("Retriever.PostRecipeResponse", retrieve.RecipeResponse{
		RecipeId: recipe_id,
		UserId: 1,
		GroupId: 1,
		// Yes if they accepted, no if they declined.
		Response: response,
	}, &success)
	
	w.Write( []byte("Success!") )
	le.Update(log.STATUS_COMPLETE, "", nil)
}

func main() {
	flag.Parse()
	le := log.New("frontend", nil)

	// Data requests (API calls)
//	http.HandleFunc("/api/ingredients", list_ingredients)
//	http.HandleFunc("/api/response", submit_response)
//	http.HandleFunc("/api/best", best_recipes)

	// Supported API calls
	// Specification at https://github.com/luke-segars/dinder-docs
	http.HandleFunc("/user", router.User)
//	http.HandleFunc("/user/login", route_user_login)
//	http.HandleFunc("/groups", route_groups)
//	http.HandleFunc("/group/add_user", route_group_add_user)
//	http.HandleFunc("/meal", route_meal)
//	http.HandleFunc("/meal/status", route_meal_status)
//	http.HandleFunc("/recipes", route_meal_status)
//	http.HandleFunc("/vote", route_vote)


	// Standard web server HTTP requests
	http.HandleFunc("/rate", rate_index_handler)
	// Serve any files in static/ directly from the filesystem.
	http.HandleFunc("/rate/static/", func(w http.ResponseWriter, r *http.Request) {
		le := log.New("web_request", log.Fields{
			"handler": "<inline>",
			"path": r.URL.Path[1:],
		})

		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
		le.Update(log.STATUS_COMPLETE, "", nil)
	})
	// No-op handler for favicon.ico, since it'll otherwise generate an extra call to index.
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	// Serve any files in static/ directly from the filesystem.
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		le := log.New("web_request", log.Fields{
			"handler": "<inline>",
			"path": r.URL.Path[1:],
		})
		
		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
		le.Update(log.STATUS_COMPLETE, "", nil)
	})

	le.Update(log.STATUS_OK, fmt.Sprintf("Awaiting requests on port %d", conf.Frontend.Port), nil)
	le.Update(log.STATUS_FATAL, "Couldn't listen on port 8088:" + http.ListenAndServe( fmt.Sprintf(":%d", conf.Frontend.Port), nil ).Error(), nil)
}
