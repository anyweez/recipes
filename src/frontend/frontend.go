package main

import (
	"flag"
	"log"
	"net/http"
)

var INGREDIENTS = flag.String("ingredients", "m/0ggm5yy", "The ingredients we should search for.")
var RETRIEVER = flag.String("retriever", "127.0.0.1:14501", "")
var MONGO = flag.String("mongo", "127.0.0.1:27017", "")

func main() {
	flag.Parse()

	http.HandleFunc("/api/ingredients", list_ingredients)
	http.HandleFunc("/api/recipes", find_recipes)
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
