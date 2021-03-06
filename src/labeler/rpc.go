package main

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/gedex/inflector"
	"log"
	proto "proto"
	"strings"
)

type Labeler int

func (l *Labeler) GetIngredient(la *LabelerArgs, reply *proto.Ingredient) error {
	log.Println(fmt.Sprintf("RPC REQUEST: [%s]", la.String()))

	// Strip out punctuation, etc (currently only comma)
	la.IngredientString = strings.Replace(la.IngredientString, ",", "", -1)
	la.IngredientString = strings.Replace(la.IngredientString, "(", "", -1)
	la.IngredientString = strings.Replace(la.IngredientString, ")", "", -1)

	reply.Name = gproto.String(la.IngredientString)

	// Lowercase the string (all strings in mapping are lowercase for case insensitivity)
	la.IngredientString = strings.ToLower(la.IngredientString)

	// Break the phrase up into tokens.
	tokens := strings.Split(la.IngredientString, " ")
	found := false

	// Iterate through all possible substrings
	for length := 1; length <= len(tokens); length++ {
		for start := 0; start+length <= len(tokens); start++ {
			log.Println(strings.Join(tokens[start:start+length], " "))
			substr := inflector.Singularize(strings.Join(tokens[start:start+length], " "))
			val, exists := IngredientMap[substr]

			// If the key exists, save it.
			if exists {
				reply.Ingrids = append(reply.Ingrids, val)
				found = true
			}
		}
	}

	// Extract all of the modifiers describing how the food should be prepared.
	// TODO: move this to the recipe parser; doesn't need to be done online.
	for _, token := range tokens {
		switch token {
		case "chopped":
			reply.Modifiers = append(reply.Modifiers, proto.Ingredient_CHOPPED)
			break
		case "grated":
			reply.Modifiers = append(reply.Modifiers, proto.Ingredient_GRATED)
			break
		case "shredded":
			reply.Modifiers = append(reply.Modifiers, proto.Ingredient_SHREDDED)
			break
		case "sliced":
			reply.Modifiers = append(reply.Modifiers, proto.Ingredient_SLICED)
			break
		case "peeled":
			reply.Modifiers = append(reply.Modifiers, proto.Ingredient_SLICED)
			break
		default:
			break
		}
	}

	log.Println(fmt.Sprintf("RPC REQUEST: [%s] found? %b", la.String(), found))
	// No errors
	return nil
}
