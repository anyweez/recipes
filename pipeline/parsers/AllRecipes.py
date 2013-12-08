import bs4
import lib.common as common

def strip_non_ascii(string):
	return ''.join(c for c in string if ord(c) < 128)

tree = None

## Not sure how I like this as a design. It works...alternative would be
## to have some sort of state object returned. The only concern I'm aware
## of would be multi-threading, but each process will only be executing a
## single codepath at one time.
def use(doc):
	global tree
	tree = bs4.BeautifulSoup(doc)
	
def clear():
	tree = None

## Parse the title out.
def title():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	try:
		return strip_non_ascii(tree.find(id='itemTitle').text)
	except AttributeError:
		return "<no name>"

## Parse the prep time.
def prep_time():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	try:
		time_str = ' '.join(tree.find(id='liPrep').text.split()[1:])	
		return common.get_time(time_str.strip())
	except AttributeError:
		return 0
			
## Parse the cook time.
def cook_time():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	
	try:
		time_str = ' '.join(tree.find(id='liCook').text.split()[1:])
		return common.get_time(time_str.strip())
	except AttributeError:
		return 0

## Parse the ready time.
def ready_time():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	
	try:
		time_str = ' '.join(tree.find(id='liTotal').text.split()[2:])
		return common.get_time(time_str.strip())
	except AttributeError:
		return 0

## Parse the number of servings.
def servings():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	
	try:
		return common.get_servings(tree.find(id='lblYield').text)
	# TODO: This is a subpar way to handle this condition. Not sure what
	#   the right way is yet though...
	except AttributeError:
		return 0

## Parse the ingredients list.
def ingredients():
	if tree is None:
		raise Exception("parsing.set() a doc first")
	
	try:
		ingr = []
		for leaf in tree.find_all(id='liIngredient'):
			cleaned = strip_non_ascii(leaf.text.strip().replace('\n',' '))

			if len(cleaned) > 0:
				ingr.append( str(cleaned) )
			
		return ingr
	except AttributeError:
		return []
