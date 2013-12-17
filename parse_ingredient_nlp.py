import nltk, sys

class Segments(object):
	QUANTITY, INGREDIENT = range(1, 3)

## Converts a string into a set of features.
## TODO: Add new features.
def features(string):
	tokens = nltk.word_tokenize(string)
#	tags = nltk.pos_tag(tokens)
	
	feat = {}
	feat['num_numbers'] = len([ l for l in string if l.isdigit() ])
	feat['num_tokens'] = len(tokens)
	feat['unit_token'] = tokens[1] if len(tokens) > 1 else ''
	
	for token in tokens:
		feat['word_%s' % (token.lower())] = True

#	tags_dict = {}
#	for token, tag in tags:
#		if tags_dict.has_key(tag):
#			tags_dict[tag] += 1
#		else:
#			tags_dict[tag] = 1
#			
#	for tag in tags_dict.keys():
#		feat['tag_%s' % tag] = tags_dict[tag]

	return feat

def train(trn):
	print 'Training on %d examples...' % len(trn)
	quantities = [ (features(t[0]), Segments.QUANTITY) for t in trn ]
	ingredients = [ (features(t[1]), Segments.INGREDIENT) for t in trn ]

	return nltk.NaiveBayesClassifier.train(quantities + ingredients)

def classify_as(classifier, first, second):
	first_class = classifier.classify( features(first) )
	first_prob = classifier.prob_classify( features(first) ).prob(first_class)
	
	second_class = classifier.classify( features(second) )
	second_prob = classifier.prob_classify( features(second) ).prob(second_class)

#	print '  (%d, %.3f, %d, %.3f)' % (first_class, first_prob, second_class, second_prob)
#	print '  (%r, %.3f)' % (first_class == Segments.QUANTITY and second_class == Segments.INGREDIENT, (first_prob + second_prob) / 2)
	if first_class == Segments.QUANTITY and second_class == Segments.INGREDIENT:
		return (first_prob + second_prob) / 2
	else:
		return 0.0

dataset = []

with open('ingredients.labeled.txt') as fp:
	lines = fp.readlines()
	for line in lines:
		first, second = line.split('\t')
		
		if len( first.strip() ) > 0:
			dataset.append( (first.strip(), second.strip()) )
		else:
			dataset.append( ('', second.strip()) )
	
# Training set should be 75% of the available 
training_size = int( len(dataset) * .75 )
classifier = train(dataset[:training_size])

num_to_run = 100
total_count = 0
correct_count = 0
for ingr in dataset[training_size:]:
	full_string = ' '.join(ingr)
	
	print full_string
	# Partition into two parts in every possible way
	tokens = nltk.word_tokenize(full_string)

	max_score = 0.0
	max_i = -1
	# Loop through all but the [full_string | ''] case, since we're assuming
	# the ingredient is always mentioned.
	for i in xrange( len(tokens) ):
#		print 'Option #%d: %s | %s' % (i + 1, ' '.join(tokens[:i]), ' '.join(tokens[i:]))
		score = classify_as(classifier, ' '.join(tokens[:i]), ' '.join(tokens[i:]))

		if score > max_score:
			max_i = i
			max_score = score
	
	print 'Full str: %s' % full_string
	print '  Quantity: %s [%r]' % (' '.join(tokens[:max_i]), ' '.join(tokens[:max_i]).replace(' ', '') == ingr[0].replace(' ', '')) 
	print '  Ingredient: %s [%r]' % (' '.join(tokens[max_i:]), ' '.join(tokens[max_i:]).replace(' ', '') == ingr[1].replace(' ', ''))
	print '  Confidence: %.2f' % score
	print ''
	
	if ' '.join(tokens[:max_i]).replace(' ', '') == ingr[0].replace(' ', '') and \
	   ' '.join(tokens[max_i:]).replace(' ', '') == ingr[1].replace(' ', ''):
		correct_count += 1
	total_count += 1
	
	if num_to_run == total_count:
		print correct_count
		print total_count
		print 'Accuracy: %f' % (float(correct_count) / total_count)
		sys.exit(0)

print correct_count
print total_count
print 'Accuracy: %f' % (float(correct_count) / total_count)
classifier.show_most_informative_features(10)
