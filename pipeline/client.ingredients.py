import juggle.lib.juggle as jugglelib
import proto.Recipes_pb2 as proto
import parsers.AllRecipes as docparser
import recipes
import os, os.path, logging, sys
import datetime as dt

## This script starts a client that batch processes a set of HTML files and generates
## structured Recipes out of them. It depends on having a running ingredients service
## that it can communicate with.
##
## A path to a directory containing recipes to be parsed should be passed as the first
## parameter.
##
## NOTE: recipe ID's are not stable across graphs.

# Configure logging.
logging.basicConfig()
logging.getLogger().setLevel(logging.INFO)

# Directory where all recipe HTML files can be found.
RECIPE_DIRECTORY = sys.argv[1]

# Sending a string, so this is a passthrough.
def encode(obj):
	return obj

# Will receive an Ingredient proto back, need to decode.
def decode(string):
	ingr = proto.Ingredient()
	ingr.ParseFromString(string)
	
	return recipes.Ingredient(ingr)

if __name__ == '__main__':
	book = proto.RecipeBook()
	
	# Set up the ingredient service to parse ingredient strings.
	ingredient_service = jugglelib.ServiceAPI('localhost', 19001)	
	ingredient_service.encoder(encode)
	ingredient_service.decoder(decode)
	
	recipe_count = len([f for f in os.listdir(RECIPE_DIRECTORY) if os.path.isfile('/'.join([RECIPE_DIRECTORY, f]))])
	count = 0
	
	for fn in [f for f in os.listdir(RECIPE_DIRECTORY) if os.path.isfile('/'.join([RECIPE_DIRECTORY, f]))]:
#		logging.info( 'Reading %s' % fn )
		with open("%s/%s" % (RECIPE_DIRECTORY, fn)) as fp: 
			doc = fp.read()
	
		logging.info( 'Parsing %s' % fn )
		docparser.use(doc)

		title = docparser.title()
		prep_time = docparser.prep_time()
		cook_time = docparser.cook_time()
		ready_time = docparser.ready_time()
		servings = docparser.servings()
		ingrs = docparser.ingredients()

		docparser.clear()

		ingredients = []
		for ingr in ingrs:
			ingredients.append( ingredient_service.query(ingr) )
		
		recipe = recipes.Recipe(
			id=count,
			title=title,
			prep_time=prep_time,
			cook_time=cook_time,
			ready_time=ready_time,
			serving_size=servings,
			ingredients=ingredients
		)
			
		# Recommended technique for adding a message to a repeated field.
		# Seems weird but it works.
		book.recipes.extend([recipe.data])

		count += 1
		print "%d / %d completed [ %.1f %% ]" % (count, recipe_count, (float(count) / float(recipe_count)) * 100)

	outfile_name = 'recipes-%s.bin' % dt.datetime.now().strftime('%y%m%d')

	print book.recipes[10].id
	print 'Writing %d recipes to %s...' % (len(book.recipes), outfile_name)
	with open(outfile_name, 'wb') as fp:
		fp.write(book.SerializeToString())
