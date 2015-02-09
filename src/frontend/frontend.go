package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	//	"io/ioutil"
	"lib/config"
	log "logging"
	"net/http"
	//	"net/rpc"
	//	retrieve "retrieve"
	router "frontend/router"
	//	"strconv"
)

var conf config.RecipesConfig
var CONFIG_LOCATION = flag.String("config", "recipes.conf", "The path to the configuration file.")

func init() {
	flag.Parse()
	c, err := config.New(*CONFIG_LOCATION)
	conf = c

	if err != nil {
		le := log.New("init", nil)
		le.Update(log.STATUS_FATAL, err.Error(), nil)
	}
}

// Fetch the index page for the "rate" URL path.
/*
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
*/

func main() {
	flag.Parse()
	le := log.New("frontend", nil)

	// Data requests (API calls)
	//	http.HandleFunc("/api/ingredients", list_ingredients)
	//	http.HandleFunc("/api/response", submit_response)
	//	http.HandleFunc("/api/best", best_recipes)

	r := mux.NewRouter()
	// Supported API calls
	// Specification at https://github.com/luke-segars/dinder-docs
	r.HandleFunc("/api/users", router.User)
	r.HandleFunc("/api/users/me", router.UserMe)
	r.HandleFunc("/api/users/login", router.Login)
	r.HandleFunc("/api/groups", router.Groups)
	r.HandleFunc("/api/groups/{group_id}/user", router.GroupsIdUser)
	r.HandleFunc("/api/meals/today", router.MealsToday)
	r.HandleFunc("/api/meals/today/status", router.MealsTodayStatus)
	r.HandleFunc("/api/meals/vote", router.MealsVote)
	r.HandleFunc("/api/recipes", router.Recipes)

	// Standard web server HTTP requests
	//	r.HandleFunc("/rate", rate_index_handler)
	// Serve any files in static/ directly from the filesystem.
	//	r.HandleFunc("/rate/static/", func(w http.ResponseWriter, r *http.Request) {
	//		le := log.New("web_request", log.Fields{
	//			"handler": "<inline>",
	//			"path": r.URL.Path[1:],
	//		})
	//
	//		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
	//		le.Update(log.STATUS_COMPLETE, "", nil)
	//	})
	// No-op handler for favicon.ico, since it'll otherwise generate an extra call to index.
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	// Serve any files in static/ directly from the filesystem.
	r.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		le := log.New("web_request", log.Fields{
			"handler": "<inline>",
			"path":    r.URL.Path[1:],
		})

		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
		le.Update(log.STATUS_COMPLETE, "", nil)
	})

	http.Handle("/", r)
	le.Update(log.STATUS_OK, fmt.Sprintf("Awaiting requests on port %d", conf.Frontend.Port), nil)

	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Frontend.Port), context.ClearHandler(http.DefaultServeMux))
	le.Update(log.STATUS_FATAL, "Couldn't listen on port 8088:"+err.Error(), nil)
}
