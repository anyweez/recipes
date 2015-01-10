package main

import (
	"bytes"
	"flag"
	"fmt"
	gproto "code.google.com/p/goprotobuf/proto"
	html "golang.org/x/net/html"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"lib/config"
	"lib/recipes"
	"log"
	"os"
	proto "proto"
	"math/rand"
	"strings"
	"time"
)

var HTML_FILES = flag.String("files", "", "The HTML file that should be parsed.")
var LABELER = flag.String("labeler", "127.0.0.1:14500", "The network location of a labeler RPC service.")
var MONGO_ADDR = flag.String("mongo", "localhost:27017", "The address for the mongo server.")
var OUTPUT_QUADS = flag.String("out", "graph.nq", "The file where the quads file should be output.")
var SAMPLE_SIZE = flag.Int("samples", 100, "The number of samples to run (only valid for `extract sample`)")

type PageRecord struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page []byte
	// The HTML content of the page.
	Content []byte
}

type GraphComponent struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Label     string `json:"label"`
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
func writeRecipe(recipe proto.Recipe, out *os.File, session *mgo.Session, conf config.RecipesConfig) {
	out.WriteString(fmt.Sprintf("<%s> <named> \"%s\" .\n", *recipe.Id, *recipe.Name))
	// Create records for each ingredient ID linking to the recipe.
	for _, ingr := range recipe.Ingredients {
		// Iterate through each mid.
		for _, iid := range ingr.Ingrids {
			out.WriteString(fmt.Sprintf("<%s> <contains> <%s> .\n", *recipe.Id, iid))
		}
	}

	// Record the structured data to Mongo.
	c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
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
	valid_modes := []string{"ingredients", "sample", "recipes"}

	// Check to ensure that a mode has been specified, and that that mode is valid.
	if len(os.Args) < 2 || !validMode(os.Args[1], valid_modes) {
		log.Fatal(fmt.Sprintf("You must specify a valid mode: [%s]", strings.Join(valid_modes, ",")))
	}
	mode := os.Args[1]
	// Load the configuration.
	conf := config.New("recipes.conf")

	switch mode {
	/**
	 * Extracts ingredients from a Freebase triples file and updates MongoDB to include
	 * all important (structured) information.
	 */
	case "ingredients":
		// Extract ingredients from the Freebase database identified in the configuration.
		ingr := ExtractIngredients(conf)
		log.Println(fmt.Sprintf("%d ingredients read in.", len(ingr)))
		// Update MongoDB.
		UpdateIngredients(conf, ingr)
		break
	case "sample":
		session, err := mgo.Dial(conf.Mongo.ConnectionString())
		if err != nil {
			log.Fatal("Cannot connect to Mongo instance: " + err.Error())
		}
		defer session.Close()

		c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RecipeCollection)
		numRecords, _ := c.Count()
		rand.Seed(time.Now().Unix())
		skip := rand.Int() % numRecords
		result := proto.Recipe{}
		
		c.Find(nil).Skip(skip).One(&result)
		recipes.DebugPrint(result)
/*
		session, err := mgo.Dial(conf.Mongo.ConnectionString())
		if err != nil {
			log.Fatal("Cannot connect to Mongo instance: " + err.Error())
		}

		defer session.Close()

		c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RawCollection)
		
		numRecords, _ := c.Count()
		var result PageRecord
		found_ingredients := 0.0
		total_ingredients := 0.0
		sample_misses := make([]string, 0)
		
		for i := 0; i < *SAMPLE_SIZE; i++ {
			fmt.Print( fmt.Sprintf("Evaluating sample... [%d / %d]\r", i+1, *SAMPLE_SIZE) )
			// Retrieve a random record.
			skipCount := rand.Int() % (numRecords - 1)
			c.Find(nil).Skip(skipCount).One(&result)
			
			// Parse the recipe and count how many ingredients were identified.
			recipe := parse(result.Content)
			
			for _, ingr := range recipe.Ingredients {
				// Record counts for stats reporting later.
				if len(ingr.Ingrids) > 0 {
					found_ingredients++
				} else {
					if rand.Int() % 100 > 10 {
						sample_misses = append(sample_misses, *ingr.Name)
					}
				}
				total_ingredients++
			}
		}
		
		if total_ingredients > 0 {
			fmt.Println( fmt.Sprintf("Ingredient recall rate: %f", found_ingredients / total_ingredients) )
		} else {
			fmt.Println("No ingredients available in parsed recipes.")
		}
*/
		break
	
	/**
	 * Parse raw HTML content and extract structured recipes. Both input and output are
	 * expected to be in MongoDB.
	 */
	case "recipes":
		output, err := os.Create(*OUTPUT_QUADS)
		if err != nil {
			log.Fatal("Couldn't open output file: " + *OUTPUT_QUADS)
		}
		
		defer output.Close()

		log.Println("Connecting to MongoDB instance at " + conf.Mongo.ConnectionString())
		session, err := mgo.Dial(conf.Mongo.ConnectionString())
		log.Println("Reading from MongoDB...")
		if err != nil {
			log.Fatal("Cannot connect to Mongo instance: " + err.Error())
		}
		defer session.Close()

		c := session.DB(conf.Mongo.DatabaseName).C(conf.Mongo.RawCollection)

		var result PageRecord
		iter := c.Find(nil).Iter()

		i := 0
		for iter.Next(&result) {
			recipe := parse(result.Content)
			recipe.SourceUrl = gproto.String( string(result.Page) )
			fmt.Println(fmt.Sprintf("%d. %s (%d min prep, %d min cook, %d min ready)",
				i+1,
				*recipe.Name,
				*recipe.Time.Prep,
				*recipe.Time.Cook,
				*recipe.Time.Ready))

			for _, ingr := range recipe.Ingredients {
				fmt.Println(fmt.Sprintf("  - %s (%s)", *ingr.Name, strings.Join(ingr.Ingrids, ", ")))
			}	

			writeRecipe(recipe, output, session, conf)
			i += 1
		}
		
		break
	}
}
