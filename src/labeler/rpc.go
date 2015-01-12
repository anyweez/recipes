package main

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"log"
	proto "proto"
	"strings"
)

type Labeler int

func (l *Labeler) GetIngredient(la *LabelerArgs, reply *proto.Ingredient) error {
	log.Println(fmt.Sprintf("RPC REQUEST: [%s]", la.String()))

	// Strip out punctuation, etc (currently only comma)
	la.IngredientString = strings.Replace(la.IngredientString, ",", " ", -1)
	la.IngredientString = strings.Replace(la.IngredientString, "(", " ", -1)
	la.IngredientString = strings.Replace(la.IngredientString, ")", " ", -1)

	reply.Name = gproto.String(la.IngredientString)

	// Break the phrase up into tokens.
	tokens := strings.Split(la.IngredientString, " ")

	// Iterate through all possible substrings
	for length := 1; length <= len(tokens); length++ {
		for start := 0; start+length < len(tokens); start++ {
			substr := strings.Join(tokens[start:start+length], " ")
			val, exists := IngredientMap[substr]

			// If the key exists, save it.
			if exists {
				reply.Ingrids = append(reply.Ingrids, val)
			}
		}
	}

	// Extract all of the modifiers describing how the food should be prepared.
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

	// No errors
	return nil
}
