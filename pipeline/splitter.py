import nltk
import recipes.proto.Recipes_pb2 as proto

class Segments(object):
	QUANTITY, INGREDIENT = range(1, 3)

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

## IngredientSplitter's convert an ingredient string like '3 tablespoon of chocolate sauce' into
## an AMOUNT component and a UNIT component using natural language processing.
class IngredientSplitter(object):
		def __init__(self):
			self.classifier = None

		## Converts a string into a set of features.
		## TODO: Add new features.
		def _features(self, string):
			tokens = nltk.word_tokenize(string)
	
			feat = {}
			feat['num_numbers'] = len([ l for l in string if l.isdigit() ])
			feat['num_tokens'] = len(tokens)
			feat['unit_token'] = tokens[1] if len(tokens) > 1 else ''
			# Are the parenthesis balanced in the segment? There are a number of error cases where
			# strings like 1 can (12 oz) is split in the middle of the parenthesis.
			feat['balanced_parens'] = len([l for l in string if l in ('(', ')')]) in (0, 2)
			# There's a fairly constrained vocabulary at the end of the QUANTITY segment, i.e.
			# cup, teaspoon, package, etc.
			feat['last_token'] = tokens[-1] if len(tokens) > 0 else ''
			feat['last_token_is_terminal'] = tokens[-1] in normalized_terms.keys() if len(tokens) > 0 else False
			feat['has_normalized_units'] = len([t for t in tokens if t in normalized_terms.keys()]) > 0

			for token in tokens:
				feat['word_%s' % (token.lower())] = True

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
