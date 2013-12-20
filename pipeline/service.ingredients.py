import juggle.lib.juggle as jugglelib
import proto.Recipes_pb2 as proto
import logging, re, sys, nltk
import lib.common as common
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

# Base units:
#   VOLUME: ounces
#   MASS: ounces
#   COUNT: (none)

# A map containing all ingredients linked to their ingrids.
ingredient_list = {}

class Segments(object):
        QUANTITY, INGREDIENT = range(1, 3)

## IngredientSplitter's convert an ingredient string like '3 tablespoon of chocolate sauce' into
## an AMOUNT component and a UNIT component using natural language processing.
class IngredientSplitter(object):
	def __init__(self):
		self.classifier = None

	## Converts a string into a set of features.
	## TODO: Add new features.
	def _features(self, string):
	        tokens = nltk.word_tokenize(string)
	#       tags = nltk.pos_tag(tokens)

	        feat = {}
	        feat['num_numbers'] = len([ l for l in string if l.isdigit() ])
	        feat['num_tokens'] = len(tokens)
	        feat['unit_token'] = tokens[1] if len(tokens) > 1 else ''

	        for token in tokens:
	                feat['word_%s' % (token.lower())] = True

	#       tags_dict = {}
	#       for token, tag in tags:
	#               if tags_dict.has_key(tag):
	#                       tags_dict[tag] += 1
	#               else:
	#                       tags_dict[tag] = 1
#                       
	#       for tag in tags_dict.keys():
	#               feat['tag_%s' % tag] = tags_dict[tag]

        	return feat

	def train(self, trn):
        	print 'Training on %d examples...' % len(trn)
        	quantities = [ (self._features(t[0]), Segments.QUANTITY) for t in trn ]
        	ingredients = [ (self._features(t[1]), Segments.INGREDIENT) for t in trn ]

        	self.classifier = nltk.NaiveBayesClassifier.train(quantities + ingredients)

	def _classify_as(self, first, second):
        	first_class = self.classifier.classify( self._features(first) )
        	first_prob = self.classifier.prob_classify( self._features(first) ).prob(first_class)

       		second_class = self.classifier.classify( self._features(second) )
        	second_prob = self.classifier.prob_classify( self._features(second) ).prob(second_class)

        	if first_class == Segments.QUANTITY and second_class == Segments.INGREDIENT:
                	return (first_prob + second_prob) / 2
        	else:
                	return 0.0

	def classify(self, ingredient_string):
		tokens = nltk.word_tokenize(ingredient_string)

		max_score = -1.0
		for i in xrange( len(tokens) ):
			score = self._classify_as(' '.join(tokens[:i]), ' '.join(tokens[i:]))

                	if score > max_score:
                        	max_i = i
                        	max_score = score

		return (' '.join(tokens[:max_i]), ' '.join(tokens[max_i:]))

#        print 'Full str: %s' % full_string
#        print '  Quantity: %s [%r]' % (' '.join(tokens[:max_i]), ' '.join(tokens[:max_i]).replace(' ', '') == ingr[0].replace(' ', ''))
#        print '  Ingredient: %s [%r]' % (' '.join(tokens[max_i:]), ' '.join(tokens[max_i:]).replace(' ', '') == ingr[1].replace(' ', ''))
#        print '  Confidence: %.2f' % score
#        print ''
splitter = IngredientSplitter()

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
		
		return ( fraction_expression.parseString(quantity_regex.group(1)), quantity_str[quantity_regex.end(1):] )
	else:
		return ( 0.0, quantity_str )

## 
def parse(qry):
	print 'starting "%s"' % qry
	
	with open('ingredient_strings.txt', 'a') as fp:
		fp.write('%s\n' % qry)
	
	# General approach: list of terms that are "units" (cup, can, etc)
	# Everything before the unit is the quantity, everything after is
	# the ingredient name.
	ingr = proto.Ingredient()

	quantity, ingredient = splitter.classify( qry )
	ingr.ingrids.extend( find_ingrids(ingredient) )

	amnt, unit = split_amount_and_unit(quantity)
	print '< %s > < %s > < %s >' % (amnt, unit, ingredient)
	
#	ingr.quantity.original_amount, ingr.quantity.original_unit_string = split_amount_and_unit( quantity )
#	ingr.quantity.original_type = normalize_type( ingr.quantity.original_unit_string )

#	ingr.quantity.normalized_amount = normalize_quantity(ingr.quantity.original_amount, ingr.quantity.original_type, ingr.quantity.original_unit_string)
	
#	ingr.amount = quantity
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

		splitter.train( training )
	
	ingredients = jugglelib.Service('ingredients', 19001)
	ingredients.handler(parse)
	
	ingredients.encoder(encode)
	ingredients.decoder(decode)
	
	ingredients.start()
