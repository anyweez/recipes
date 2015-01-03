package main

import (
	html "golang.org/x/net/html"
	"bytes"
	"flag"
//	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
//	"net/http"
	"io/ioutil"
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

func main() {
	flag.Parse()

	output, _ := os.Create(*OUTPUT_QUADS)
	defer output.Close()
	
	/**
	 * Read from disk.
	 */
	if len(*HTML_FILES) > 0 {
		log.Println("Reading files from disk.")
		
		files, err := ioutil.ReadDir(*HTML_FILES)
		if err != nil {
			log.Fatal("Error reading files:" + err.Error())
		}
	
		for i, file := range files {
			// If this file is an HTML file, 
			if strings.HasSuffix(file.Name(), ".html") {
				// Open the file and pass the data to parse()
				data, err := ioutil.ReadFile(*HTML_FILES + "/" + file.Name())
				if err != nil {
					log.Fatal("Couldn't open file:" + err.Error())
				}

				recipe := parse(data)
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
			}
		}
	/**
	 * Connect to Mongo instance.
	 */
	} else {
		log.Println("Reading from MongoDB instance.")
		
		session, err := mgo.Dial(*MONGO_ADDR)
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
	}
	
}
