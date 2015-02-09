package main

import (
	gproto "code.google.com/p/goprotobuf/proto"
	"fmt"
	html "golang.org/x/net/html"
	labeler "labeler"
	"log"
	"net/rpc"
	proto "proto"
	"strconv"
	"strings"
)

var nextId uint32

/**
 * Initialize nextId to zero.
 */
func init() {
	nextId = 0
}

func isTitleToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == "itemTitle" {
			return true
		}
	}

	return false
}

func isPrepTimeToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == "timePrep" && token.Data == "time" {
			return true
		}
	}

	return false
}

func isCookTimeToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == "timeCook" && token.Data == "time" {
			return true
		}
	}

	return false
}

func isReadyTimeToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == "timeTotal" && token.Data == "time" {
			return true
		}
	}

	return false
}

func isIngredientToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "itemprop" && attr.Val == "ingredients" {
			return true
		}
	}

	return false
}

func isImageToken(token html.Token) bool {
	for _, attr := range token.Attr {
		if attr.Key == "id" && attr.Val == "imgPhoto" {
			return true
		}
	}

	return false
}

func _parser(tk *html.Tokenizer) proto.Recipe {
	client, err := rpc.DialHTTP("tcp", *LABELER)
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't connect to labeler at %s [%s]", *LABELER, err.Error()))
	}
	defer client.Close()

	recipe := proto.Recipe{}
	recipe.Id = gproto.String(fmt.Sprintf("/r/%d", nextId))
	nextId += 1

	recipe.Time = &proto.Recipe_Time{
		Prep:  gproto.Uint32(0),
		Cook:  gproto.Uint32(0),
		Ready: gproto.Uint32(0),
	}

	next := tk.Next()
	for next != html.ErrorToken {
		tok := tk.Token()

		// Parser
		if tok.Type == html.StartTagToken || tok.Type == html.SelfClosingTagToken {

			/**
			 * Title tokens contain the title of the recipe. The tokken
			 * after the title token is what contains the title itself.
			 */
			if isTitleToken(tok) {
				tk.Next()
				recipe.Name = gproto.String(_getTitle(tk.Token()))
				/**
				 * The PrepTimeToken contains the number of hours and minutes
				 * required to prepare the dish. It is extracted from the <time>
				 * element.
				 */
			} else if isPrepTimeToken(tok) {
				recipe.Time.Prep = gproto.Uint32(_getTime(tok))
			} else if isCookTimeToken(tok) {
				recipe.Time.Cook = gproto.Uint32(_getTime(tok))
			} else if isReadyTimeToken(tok) {
				recipe.Time.Ready = gproto.Uint32(_getTime(tok))
			} else if isIngredientToken(tok) {
				hasQuantity := true
				quantity := ""
				ingrName := ""

				// Parse out the ingredient name and the quantity.
				for len(ingrName) == 0 && (hasQuantity || len(quantity) == 0) {
					tk.Next()
					tok = tk.Token()
					for _, attr := range tok.Attr {
						if attr.Key == "class" && attr.Val == "ingredient-amount" {
							tk.Next()
							quantity = tk.Token().String()
						}

						if attr.Key == "class" && attr.Val == "ingredient-name" {
							tk.Next()
							ingrName = tk.Token().String()

							// If we get the ingredient name with no quantity value,
							// mark that the quality value won't be coming.
							if len(quantity) == 0 {
								hasQuantity = false
							}
						}

						// On a rare occasion (< .1%) recipes will have a special type
						// of heading that seems to be able to be treated the same as
						// ingredient-name. Keeping it separate for now since it seems
						// like it could be semantically different, or accidental on the
						// part of the recipe provider.
						if attr.Key == "class" && attr.Val == "ingred-heading" {
							tk.Next()
							ingrName = tk.Token().String()

							if len(quantity) == 0 {
								hasQuantity = false
							}
						}
					}
				}
				recipe.Ingredients = append(recipe.Ingredients, _getIngredient(ingrName, quantity, client))
			} else if isImageToken(tok) {
				recipe.ImageUrls = append(recipe.ImageUrls, _getImage(tok))
			}
		}

		// Get next token.
		next = tk.Next()
	}

	return recipe
}

func _getTitle(token html.Token) string {
	return html.UnescapeString(token.String())
}

func _getTime(token html.Token) uint32 {
	for _, attr := range token.Attr {
		if attr.Key == "datetime" {
			minutes := 0
			str := attr.Val[2:len(attr.Val)]

			// Parse the # of hours and add to accumulator.
			hr_index := strings.Index(str, "H")
			if hr_index >= 0 {
				hours_delta, _ := strconv.Atoi(str[0:hr_index])
				minutes += hours_delta * 60
				str = str[hr_index+1 : len(str)]
			}

			// Parse the # of minutes and add to accumulator.
			min_index := strings.Index(str, "M")
			if min_index >= 0 {
				minutes_delta, _ := strconv.Atoi(str[0:min_index])
				minutes += minutes_delta
				str = str[min_index+1 : len(str)]
			}

			return uint32(minutes)
		}
	}

	return 0
}

func _getIngredient(in string, quantity string, client *rpc.Client) *proto.Ingredient {
	ingr := proto.Ingredient{}

	err := client.Call("Labeler.GetIngredient", labeler.LabelerArgs{
		IngredientString: html.UnescapeString(in),
	}, &ingr)

	if err != nil {
		log.Fatal("Error parsing ingredient:" + err.Error())
	}
	ingr.QuantityString = gproto.String(quantity)

	return &ingr
}

func _getImage(token html.Token) string {
	for _, attr := range token.Attr {
		if attr.Key == "src" {
			return attr.Val
		}
	}

	return ""
}
