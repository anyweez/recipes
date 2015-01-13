package main

/**
 * Labeler is a service that runs and responds to RPC's from connected
 * clients. It receives strings like "1 cup of olive oil" and returns
 * a structured representation of the string.
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"lib/config"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
)

var PORT = flag.Int("port", 14500, "The port the process should listen on.")
var INGREDIENT_DB = flag.String("db", "data/ingredients.list", "A file containing triples that describe the name and mid of each ingredient.")

// Big map that contains name => mid mappings, i.e. "onion" => "/m/0dj75"
var IngredientMap map[string]string

type LabelerArgs struct {
	IngredientString string
	QuantityString   string
}

func (l *LabelerArgs) String() string {
	return fmt.Sprintf("%s %s", l.QuantityString, l.IngredientString)
}

/**
 * This function reads in and inverts ingredients stored in MongoDB (extracted
 * with `extract ingredients` command). The mapping is then used for exact
 * string matches via the labeler rpc's.
 *
 * The mapping is case-insensitive.
 */
func loadMapping(conf config.RecipeConfig) {
	IngredientMap = make(map[string]string)

	session, err := mgo.Dial(conf.Mongo.ConnectionString())
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.IngredientCollection)

	iter := c.Find(nil).Iter()
	ingredient := proto.Ingredient{}

	// Create the name=>mid mapping.
	for iter.Next(&ingredient) {
		cleanedName := strings.TrimSpace( strings.ToLower(ingredient.Name) )
		IngredientMap[cleanedName] = ingrid
	}
}

func main() {
	flag.Parse()

	// Load the ingredient map.
	log.Println("Loading ingredient map...")
	conf := config.New("recipes.conf")
	loadMapping(conf)
	log.Println(fmt.Sprintf("Loaded %d ingredients.", len(IngredientMap)))

	// Set up the RPC HTTP interface
	labeler := new(Labeler)
	rpc.Register(labeler)
	rpc.HandleHTTP()

	// Start listening!
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))

	if err != nil {
		log.Fatal("Couldn't start listening:" + err.Error())
	}

	log.Println("Setup complete. Listening for RPC's on HTTP interface.")
	http.Serve(l, nil)
}
