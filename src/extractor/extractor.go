package main

import (
	html "golang.org/x/net/html"
	"bytes"
	"flag"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"io/ioutil"
	proto "proto"
	"strings"
)

var HTML_FILES = flag.String("files", "", "The HTML file that should be parsed.")
var LABELER = flag.String("labeler", "127.0.0.1:14500", "The network location of a labeler RPC service.")
var MONGO_ADDR = flag.String("mongo", "historian:27017", "The address for the mongo server.")

type PageRecord struct {
	Id		bson.ObjectId `bson:"_id,omitempty"`
	// The URL of the page.
	Page	[]byte
	// The HTML content of the page.
	Content	[]byte
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

func main() {
	flag.Parse()

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
					fmt.Println( fmt.Sprintf("  - %s", *ingr.Name) )
				}
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
		
		results := make([]PageRecord, 0)
		c.Find(nil).All(&results)
		
		for i, page := range results {
			recipe := parse(page.Content)
			fmt.Println( fmt.Sprintf("%d. %s (%d min prep, %d min cook, %d min ready)", 
				i+1, 
				*recipe.Name, 
				*recipe.Time.Prep,
				*recipe.Time.Cook,
				*recipe.Time.Ready) )
			
			for _, ingr := range recipe.Ingredients {
				fmt.Println( fmt.Sprintf("  - %s", *ingr.Name) )
			}
		}
	}
	
}
