package main

import (
	html "golang.org/x/net/html"
	"bytes"
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"log"
	proto "proto"
	"os"
	"strings"
)

var HTML_FILES = flag.String("files", "", "The HTML file that should be parsed.")
var LABELER = flag.String("labeler", "127.0.0.1:14500", "The network location of a labeler RPC service.")
var MONGO_ADDR = flag.String("mongo", "localhost:27017", "The address for the mongo server.")
var OUTPUT_QUADS = flag.String("out", "graph.nq", "The file where the quads file should be output.")

type PageRecord struct {
	Id		bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page	[]byte
	// The HTML content of the page.
	Content	[]byte
}

type GraphComponent struct {
	Subject		string	`json:"subject"`
	Predicate	string	`json:"predicate"`
	Object		string	`json:"object"`
	Label		string	`json:"label"`
}

/**
 * Parses an HTML file and extracts a structured recipe.
 */
func parse(data []byte) proto.Recipe {
	// Create a reader on the byte stream.
	reader := bytes.NewReader(data)

	// Create a tokenizer.
	tk := html.NewTokenizer(reader)
	return _parser(tk)
}

/**
 * Output the recipe to Cayley's HTTP endpoint. This function selects the
 * important fields from the Recipe data structure and posts them to Cayley.
 */
func writeRecipe(recipe proto.Recipe, out *os.File) {
	out.WriteString( fmt.Sprintf("<%s> <named> \"%s\" .\n", *recipe.Id, *recipe.Name) )
	// Create records for each ingredient ID linking to the recipe.
	for _, ingr := range recipe.Ingredients {
		// Iterate through each mid.
		for _, iid := range ingr.Ingrids {	
			out.WriteString( fmt.Sprintf("<%s> <contains> <%s> .\n", *recipe.Id, iid) )
		}
	}
	
	// Record the structured data to Mongo.
	session, _ := mgo.Dial(*MONGO_ADDR)
	c := session.DB("recipes").C("parsed")
	c.Insert(recipe)
}

/**
 * Checks the list of valid modes to determine whether the specified mode
 * is recognized. Return true if so, false otherwise.
 */
func validMode(target string, valid []string) bool {
	for _, val := range valid {
		if target == val {
			return true
		}
	}
	
	return false
}

func main() {
	flag.Parse()
	valid_modes := []string{ "ingredients", "recipes" }

	// Check to ensure that a mode has been specified, and that that mode is valid.
	if len(os.Args) < 2 || !validMode(os.Args[1], valid_modes) {
		log.Fatal( fmt.Sprintf("You must specify a valid mode: [%s]", strings.Join(valid_modes, ",")) )
	}
	mode := os.Args[1]
	// Load the configuration.
	conf := config.New("recipes.conf")

	switch (mode) {
		/**
		 * Extracts ingredients from a Freebase triples file and updates MongoDB to include
		 * all important (structured) information.
		 */
		case "ingredients":
			// Extract ingredients from the Freebase database identified in the configuration.
			ingr := ExtractIngredients(conf)
			log.Println( fmt.Sprintf("%d ingredients read in.", len(ingr)) )
			// Update MongoDB.
			UpdateIngredients(conf, ingr)
			break
		/**
		 * Parse raw HTML content and extract structured recipes. Both input and output are
		 * expected to be in MongoDB.
		 */
		case "recipes":
			output, _ := os.Create(*OUTPUT_QUADS)
			defer output.Close()

			log.Println("Reading from MongoDB instance.")
		
			session, err := mgo.Dial(conf.Mongo())
			if err != nil {
				log.Fatal("Cannot connect to Mongo instance: " + err.Error())
			}
		
			defer session.Close()
		
			c := session.DB("recipes").C("scraper")
		
			var result PageRecord
			iter := c.Find(nil).Iter()
		
			i := 0
			for iter.Next(&result) {
				recipe := parse(result.Content)
				fmt.Println( fmt.Sprintf("%d. %s (%d min prep, %d min cook, %d min ready)", 
					i+1, 
					*recipe.Name, 
					*recipe.Time.Prep,
					*recipe.Time.Cook,
					*recipe.Time.Ready) )
			
				for _, ingr := range recipe.Ingredients {
					fmt.Println( fmt.Sprintf("  - %s (%s)", *ingr.Name, strings.Join(ingr.Ingrids, ", ")) )
				}
			
				writeRecipe(recipe, output)
				i += 1
			}

			break
	}
}
