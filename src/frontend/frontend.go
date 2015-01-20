package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	retrieve "retrieve"
	"strconv"
)

var INGREDIENTS = flag.String("ingredients", "m/0ggm5yy", "The ingredients we should search for.")
var RETRIEVER = flag.String("retriever", "127.0.0.1:14501", "")
var MONGO = flag.String("mongo", "127.0.0.1:27017", "")

// Fetch the index page for the "rate" URL path.
func rate_index_handler(w http.ResponseWriter, r *http.Request) {
	log.Println("index requested")

	data, err := ioutil.ReadFile("html/rate/index.html")
	if err != nil {
		log.Println("rate/index.html not present!")
		http.NotFound(w, r)
	} else {
		w.Write(data)
	}
}

func submit_response(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", *RETRIEVER)
	if err != nil {
		log.Fatal("Couldn't connect to retriever: " + err.Error())
	}

	recipe_id := r.URL.Query().Get("recipe")
	response, err := strconv.ParseBool(r.URL.Query().Get("response")) // "true" or "false"
	if err != nil {
		log.Println("Invalid valid for response; must be 'true' or 'false'.")
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
}

func main() {
	flag.Parse()

	// Serve any files in static/ directly from the filesystem.
	http.HandleFunc("/rate/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL.Path[1:])
		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
	})
	// Page requests (HTML, CSS, JS, etc)
	http.HandleFunc("/rate", rate_index_handler)
	
	// Data requests (API calls)
	http.HandleFunc("/api/ingredients", list_ingredients)
	http.HandleFunc("/api/response", submit_response)
	http.HandleFunc("/api/best", best_recipes)
	// No-op handler for favicon.ico, since it'll otherwise generate an extra call to index.
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	// Serve any files in static/ directly from the filesystem.
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET", r.URL.Path[1:])
		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
	})

	log.Println("Awaiting requests...")
	log.Fatal("Couldn't listen on port 8088:", http.ListenAndServe(":8088", nil))
}
