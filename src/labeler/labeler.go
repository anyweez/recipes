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
 * This function reads in tuples for all ingredients that should be identified
 * by the labeler. It should contain one tuple per ingredient that identifies
 * the mid and the labeled name of the ingredient, i.e.
 *
 * 		/m/18df53	name	Onion
 *
 * The mapping is case-insensitive.
 */
func loadMapping(filename string) {
	IngredientMap = make(map[string]string)

	data, _ := ioutil.ReadFile(filename)
	records := strings.Split(string(data), "\n")

	for _, record := range records {
		pair := strings.Split(record, "\t")
		if len(pair) != 2 {
			log.Println("WARNING: Invalid ingredient mapping file; incomplete line doesn't include full pair:" + record)
		} else {
			cleaned := strings.TrimSpace( strings.ToLower(pair[1]) )
			IngredientMap[cleaned] = pair[0]
		}
	}
}

func main() {
	flag.Parse()

	// Load the ingredient map.
	log.Println("Loading ingredient map...")
	loadMapping(*INGREDIENT_DB)
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
