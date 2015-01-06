from py2neo import neo4j, node, rel		# neo4j support
from pymongo import MongoClient
from bson.binary import Binary
import sys, logging
import proto.Recipes_pb2 as proto

logging.basicConfig()
logger = logging.getLogger('activate')
logger.setLevel(logging.INFO)

NEO4J_IP = '10.1.1.55'
NEO4J_PORT = 7474

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

logger.info('Loading new recipes')
book = load_recipes(sys.argv[1])

if len(book.recipes) is 0:
	logger.warn("Couldn't load recipes from specified recipe book. Quitting")
	sys.exit(1)

logger.info('Connecting to neo4j')
db = neo4j.GraphDatabaseService('http://%s:%d/db/data' % (NEO4J_IP, NEO4J_PORT))
logger.info('Connecting to MongoDB')
docdb = MongoClient('mongodb://localhost:27017/')

logger.info('Clearing old graph and deleting indexes')
# Clear old state
db.clear()
#db.delete_index(neo4j.Node, 'Ingredients')
#db.delete_index(neo4j.Node, 'Recipes')

logger.info('Clearing old recipes documents')
docdb.docs.recipes.drop()
docs = docdb.docs.recipes

logger.info('Creating new indexes')
# Create new state
ingredients = db.get_or_create_index(neo4j.Node, "Ingredients")
recipes = db.get_or_create_index(neo4j.Node, "Recipes")

logger.info('Copying %d recipes into graph and document stores' % len(book.recipes))
for recipe in book.recipes:
	# TODO: add to Mongo
	docs.insert( {'_id': recipe.id, 'data': Binary(recipe.SerializeToString())} )

	# TODO: add recipe node to neo4j
	recipe_node, = db.create( node(id=recipe.id, name=recipe.name) )
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
