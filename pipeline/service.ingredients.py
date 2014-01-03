import juggle.lib.juggle as jugglelib
import proto.Recipes_pb2 as proto
import logging, re, sys, nltk
import lib.common as common
import splitter
from pyparsing import Regex, Optional

## This script starts a service that will convert strings like [1 cup tomato sauce]
## into a structured ingredient. It will also identify terms listed in freebase
## and label them accordingly.
##
## Accepts a STRING and returns an Ingredient proto.
##
## First parameter is a filename that maps ingredient id's (ingrids) to
## text tokens that should be labeled in this service.

# Configure logging.
logging.basicConfig()
logging.getLogger().setLevel(logging.INFO)

# A map containing all ingredients linked to their ingrids.
ingredient_list = {}
ingr_splitter = splitter.IngredientSplitter()

normalized_terms = {
	'cup':	proto.Ingredient.CUP,
	'cups': proto.Ingredient.CUP,
	'tablespoon': proto.Ingredient.TABLESPOON,
	'tablespoons': proto.Ingredient.TABLESPOON,
	'teaspoon': proto.Ingredient.TEASPOON,
	'teaspoons': proto.Ingredient.TEASPOON,
	'quart': proto.Ingredient.QUART,
	'quarts': proto.Ingredient.QUART,
	'gallon': proto.Ingredient.GALLON,
	'gallons': proto.Ingredient.GALLON,
	'pound': proto.Ingredient.POUND,
	'pounds': proto.Ingredient.POUND,
	'pint': proto.Ingredient.PINT,
	'pints': proto.Ingredient.PINT
}

## Identifies all known ingrids (ingredient ID's, mapped back to freebase
## entities) that are present in the ingredient name. The graph database
## will be indexed by ingrids, so if no ingrids are detected (or the
## correct ones are not detected) then it won't be possible to find this
## recipe by this ingredient.
def find_ingrids(name_str):
	# Add whitespace at the beginning and end to keep the substring algorithm
	# simple and quick.
	name_str = ' %s ' % name_str
	ingrids = [ingredient_list[key] for key in ingredient_list if ' %s ' % key in name_str or ' %ss ' % key in name_str]
	
	return ingrids

## extract_amount returns a string
def split_amount_and_unit(quantity_str):
	# This should really identify tokens that are only [0-9/]
	quantity_regex = re.search( '^([0-9 /.]+)', quantity_str )

	if quantity_regex:
	        number = Regex(r"\d+(\.\d*)?").setParseAction(lambda t: float(t[0]))
	        fraction = number('numerator') + '/' + number('denominator')
        	fraction.setParseAction(lambda t: t.numerator / t.denominator)

	        fraction_expression = fraction | number + Optional(fraction)
        	fraction_expression.setParseAction(lambda t: sum(t))
		
		return ( fraction_expression.parseString(quantity_regex.group(1))[0], quantity_str[quantity_regex.end(1):] )
	else:
		return ( 0.0, quantity_str )

## Looks for known terms that identify a standardized unit of measurement. This
## function returns a quantity type if it can be found; otherwise it returns
## the CUSTOM quantity type.
##
## Note that it is currently returning the definition of the first token
## that it comes across. This seems to be the desired behavior in practice,
## and having more than one token in a string is unlikely.
def normalize_type(unit_string):
	tokens = nltk.word_tokenize(unit_string)

	for token in tokens:
		try:
			return normalized_terms[token]
		except KeyError:
			continue
	
	return proto.Ingredient.CUSTOM

## 
def parse(qry):
	print 'starting "%s"' % qry
	
	# General approach: list of terms that are "units" (cup, can, etc)
	# Everything before the unit is the quantity, everything after is
	# the ingredient name.
	ingr = proto.Ingredient()

	quantity, ingredient = ingr_splitter.classify( qry )
	ingr.ingrids.extend( find_ingrids(ingredient) )

	amnt, unit = split_amount_and_unit(quantity)
	print '< %s > < %s > < %s >' % (amnt, unit, ingredient)
	
	ingr.quantity.original_amount, ingr.quantity.original_unit_string = split_amount_and_unit( quantity )
	ingr.quantity.original_type = normalize_type( ingr.quantity.original_unit_string )

	# TODO: normalize the amount to a standard unit
	# TODO: extract ingredient modifiers
	
	ingr.name = ingredient

	print 'completed query'
	return ingr

# Encode as an Ingredient proto.
def encode(obj):
	return obj.SerializeToString()

# Query is a string, this is a passthrough.	
def decode(string):
	return string

if __name__ == '__main__':
	print 'Reading full ingredient list...'
	with open(sys.argv[1]) as fp_ingredients:
		lines = [line.strip() for line in fp_ingredients.readlines()]
		
		for line in lines:
			try:
				ingrid, name = line.split('\t')
				ingredient_list[name.lower()] = ingrid
			except ValueError:
				print 'Warning: incomplete ingredient line: %s' % line

	print 'Training ingredient splitter...'
	with open('data/ingredients.labeled.txt') as fp_training:
		training = []
		for line in fp_training.readlines():
			if len( line.strip().split('\t') ) > 1:
				training.append( line.strip().split('\t') )
			else:
				training.append( ('', line.strip()) )

		ingr_splitter.train( training )
	
	ingredients = jugglelib.Service('ingredients', 19001)
	ingredients.handler(parse)
	
	ingredients.encoder(encode)
	ingredients.decoder(decode)
	
	ingredients.start()
