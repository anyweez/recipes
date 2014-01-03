import splitter
import sys, random

NUM_TRIALS = 20

def compare(str1, str2):
	str1 = str1.strip().replace(' ', '')
	str2 = str2.strip().replace(' ', '')

	return str1 == str2

print 'Running %d quality trials...' % NUM_TRIALS

with open(sys.argv[1]) as fp:
	lines = [tuple(line.split('\t')) for line in fp.readlines()]

scores = []
for i in xrange(NUM_TRIALS):
	ingr_splitter = splitter.IngredientSplitter()

	# Randomly choose 3/4 of lines and train
	training = random.sample( lines, int(len(lines) * .75) )

	# For remaining 1/4, try to classify and count successes
	testing = list( set(lines) - set(training) )
	ingr_splitter.train( training )
	
	# Average successes for each iteration
	successes = 0
	attempts = 0

	for test in testing:
		first, second = ingr_splitter.classify( ' '.join(test) )

		if compare(test[0], first) and compare(test[1], second):
			successes += 1
		else:
			print '%s | %s' % (first, second)

		attempts += 1

	scores.append( float(successes) / float(attempts) )
	print '[Attempt %d / %d] Score: %.2f' % (i + 1, NUM_TRIALS, float(successes) / float(attempts) )

print 'Score: %.2f over %d trials' % ( (sum(scores) / len(scores)), len(scores) )
