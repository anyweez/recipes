from py2neo import neo4j, node, rel		# neo4j support
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

# Clear old state
db.clear()
#db.delete_index(neo4j.Node, 'Ingredients')
#db.delete_index(neo4j.Node, 'Recipes')

# Create new state
ingredients = db.get_or_create_index(neo4j.Node, "Ingredients")
recipes = db.get_or_create_index(neo4j.Node, "Recipes")

book = load_recipes(sys.argv[1])

for recipe in book.recipes:
	# TODO: add to Mongo
	
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

sys.exit(0)

# Add all ingredients.
with open('ingredients.list') as fp:
	lines = [line.strip() for line in fp.readlines()]
	
	for line in lines:
		ingrid, name = line.split('\t')
		ingr_node, = db.create( node(id=ingrid, name=name) )
		# Add to the ingredients index.
		ingredients.add_if_none('id', ingrid, ingr_node)
		
		print ingredients.get('id', ingrid)

# Add all recipes.
with open('recipes.list') as fp:
	lines = [line.strip() for line in fp.readlines()]
	
	for line in lines:
		rid, name = line.split('\t')
		recipe_node, = db.create( node(id=rid, name=name) )
		# Add to the recipes index.
		recipes.add_if_none('id', rid, recipe_node)

# Add all edges.
with open('edges.list') as fp:
	lines = [line.strip() for line in fp.readlines()]
	
	for line in lines:
		ingrid, rid = line.split('\t')
		recipe, = recipes.get('id', rid)
		ingredient, = ingredients.get('id', ingrid)

		db.create( rel(recipe, 'CONTAINS', ingredient) )


sys.exit(0)

i = [ ('Celery', '/m/123'),
      ('Tofu', '/m/456')]
r = [ ('Celery + Tofu Salad', '0891', ['/m/123', '/m/456']), ]

relationships = [ ('0891', '/m/123'), ('0891', '/m/456') ]

mapping = {}

for ingr in i:
	item = node(id=ingr[1], name=ingr[0])
	db_item, = db.create(item)
	ingredients.add_if_none('id', ingr[1], db_item)

	print ingredients.get('id', ingr[1])

	mapping[ingr[1]] = db_item

#sys.exit(0)

for rec in r:
	item = node(id=rec[1], name=rec[0])
	db_item, = db.create(item)
	recipes.add_if_none('id', rec[1], db_item)
	
	print recipes.get('id', rec[1])
	
#	for id in rec[2]:
#		db.create(rel(db_item, 'CONTAINS', mapping[id]))

for src, dest in relationships:
	src_item = recipes.get('id', src)[0]
	dest_item = ingredients.get('id', dest)[0]
	db.create( rel(src_item, 'CONTAINS', dest_item) )

#print db.create(
#	node(id='/m/123', name='Celery'),
#	node(id='/m/456', name='Tofu'),
#	node(id='0891', name='Celery + Tofu Salad'),
#	rel(2, "CONTAINS", 0),
#	rel(2, "CONTAINS", 1)ate
#)
