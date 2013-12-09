from py2neo import neo4j, node, rel		# neo4j support
from pymongo import MongoClient
import sys
import proto.Recipes_pb2 as proto

## This script depends on having a RecipeBook written to a readable location
## on disk. It will load the most recent RecipeBook and add all Recipe
## objects to a document store (MongoDB) and all ingrids and their
## corresponding recipe to the graph database (neo4j).

# Load all recipes into a proto.
def load_recipes(path):
	with open(path) as fp:
		book = proto.RecipeBook()
		book.ParseFromString(fp.read())

		return book

db = neo4j.GraphDatabaseService('http://localhost:7474/db/data')
docdb = MongoClient('mongodb://localhost:27017/')
docs = docdb.recipes

# Clear old state
db.clear()
db.delete_index(neo4j.Node, 'Ingredients')
db.delete_index(neo4j.Node, 'Recipes')

# Create new state
ingredients = db.get_or_create_index(neo4j.Node, "Ingredients")
recipes = db.get_or_create_index(neo4j.Node, "Recipes")

book = load_recipes(sys.argv[1])

for recipe in book.recipes:
	# TODO: add to Mongo
	docs.insert( {'_id': recipe.id, 'data': recipe.SerializeToString()} )

	# TODO: add recipe node to neo4j
	recipe_node, = db.create( node(id=recipe.id, name=recipe.title) )
	recipes.add_if_none('id', recipe.id, recipe_node)
	
	for ingredient in recipe.ingredients:
		for ingrid in ingredient.ingrids:
			# TODO: add ingredient node if it doesn't exist
			ingredient_node = ingredients.get('id', ingrid)

			if len(ingredient_node) is 0:
				ingredient_node, = db.create( node(id=ingrid, name=ingredient.name) )
				ingredients.add_if_none('id', ingrid, ingredient_node)
			else:
				ingredient_node = ingredient_node[0]

			# TODO: add edge between ingredient and recipe
			db.create( rel(recipe_node, 'CONTAINS', ingredient_node) )
