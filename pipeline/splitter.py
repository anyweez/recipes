import nltk

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


