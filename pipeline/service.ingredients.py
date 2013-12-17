import juggle.lib.juggle as jugglelib
import proto.Recipes_pb2 as proto
import logging, re, sys
import lib.common as common

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

# Base units:
#   VOLUME: ounces
#   MASS: ounces
#   COUNT: (none)

unit_map = {
#	"can": (1, proto.Ingredient.UNKNOWN),
#	"cans": (1, proto.Ingredient.UNKNOWN),
	"cup": (8, proto.Ingredient.VOLUME),
	"cups": (8, proto.Ingredient.VOLUME),
#	"envelope": (1, proto.Ingredient.UNKNOWN),
	"oz": (1, proto.Ingredient.VOLUME),
	"ounce": (1, proto.Ingredient.VOLUME),        # always a volume?
	"ounces": (1, proto.Ingredient.VOLUME),
#	"package": (1, proto.Ingredient.UNKNOWN),
	"pound": (16, proto.Ingredient.MASS),
	"tablespoon": (.5, proto.Ingredient.VOLUME),
	"tablespoons": (.5, proto.Ingredient.VOLUME),
	"teaspoon": (.166667, proto.Ingredient.VOLUME),
	"teaspoons": (.166667, proto.Ingredient.VOLUME)
}

## Abstract units aren't linked to a particular quantity in base units.
## They should still be a part of the ingredient description.
abstract_units = [
	"package",
	"packages",
	"envelope",
	"envelopes",
	"pinch",
	"pinches",
	"bottle",
	"bottles"
]

# A map containing all ingredients linked to their ingrids.
ingredient_list = {}

## Accepts a text string and splits it into:
##   AMOUNT: numerical quantity of the ingredient
##   UNIT:   the type of unit that the amount is represented in (proto.Ingredient.*)
##   ABSTRACT_UNIT: if UNIT == UNKNOWN, this is a string describing the unknown unit
##   NAME:   the name of the ingredient
def ingr_split(text):
	tokens = [tok.replace('(', '').replace(')', '') for tok in text.split()]
	# 1) look for unit from unit_map, split on it
	units = [unit for unit in unit_map.keys() if unit in tokens]
	
	# If no units are mentioned, will be an UNKNOWN unit type and this
	# should try to figure out what the abstract unit is; falls into
	# two categories:
	#   - [1 package seasoning]
	#   - [1 egg]
	if len(units) is 0:
		amount = common.get_num(tokens[0])
		units = [u for u in abstract_units if u in tokens]
		
		if len(units) > 0:
			name = ' '.join( text.split(units[0])[1:] ).strip()
			print units[0]
			return ( amount, proto.Ingredient.UNKNOWN, units[0], name )
		else:
			return ( amount, proto.Ingredient.UNKNOWN, None, ' '.join(tokens[1:]).strip() )
			
	# If there's one unit, this is the one that we should use.
	# Two cases:
	#   - [1 pound of lean ground beef]
	#   - [2 (26 ounce) jars spaghetti sauce] 
	elif len(units) is 1:
		match = re.match(r'([0-9]+) \(([0-9./]+) ([a-zA-Z]+)\) ([a-zA-Z]+)', text)
		# Case 2
		if match:
			containers = common.get_num(match.group(1))
			size = common.get_num(match.group(2))
			unit = unit_map[match.group(3)]
			abstract_unit = match.group(4)
						
			return ( containers * size * unit[0], 
					  unit[1], 
					  None,
					  text.split(abstract_unit)[1].strip() )
		else:
			i = tokens.index(units[0])
			
			amount = ' '.join(tokens[:i])
			name = ' '.join(tokens[i+1:])
			unit = unit_map[units[0]]
			
			return ( common.get_num(amount) * unit[0], unit[1], None, name.strip() )
	else:
		raise Exception("More than one unit identified in ingredient: %s", text)

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

## 
def parse(qry):
	print 'starting "%s"' % qry
	
	with open('ingredient_strings.txt', 'a') as fp:
		fp.write('%s\n' % qry)
	
	# General approach: list of terms that are "units" (cup, can, etc)
	# Everything before the unit is the quantity, everything after is
	# the ingredient name.
	ingr = proto.Ingredient()
	
	amount, unit, abstract_unit, name = ingr_split(qry)
	ingr.ingrids.extend( find_ingrids(name) )
	
	if unit is proto.Ingredient.UNKNOWN and abstract_unit is not None:
		ingr.quantity_type = abstract_unit
	
	ingr.quantity = unit
	ingr.amount = amount
	ingr.name = name

	print 'completed query'
	return ingr

# Encode as an Ingredient proto.
def encode(obj):
	return obj.SerializeToString()

# Query is a string, this is a passthrough.	
def decode(string):
	return string

if __name__ == '__main__':
	with open(sys.argv[1]) as fp_ingredients:
		lines = [line.strip() for line in fp_ingredients.readlines()]
		
		for line in lines:
			try:
				ingrid, name = line.split('\t')
				ingredient_list[name.lower()] = ingrid
			except ValueError:
				print 'Warning: incomplete ingredient line: %s' % line
	
	ingredients = jugglelib.Service('ingredients', 19001)
	ingredients.handler(parse)
	
	ingredients.encoder(encode)
	ingredients.decoder(decode)
	
	ingredients.start()
