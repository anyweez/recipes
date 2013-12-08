import proto.Recipes_pb2 as proto

## Contains data structures used for storing, printing, and manipulating
## recipe-related data.

class Ingredient(object):
	def __init__(self, ingr):
		self.data = ingr
		
	def amount(self):
		return self.data.amount
	
	def name(self):
		return self.data.name
	
	def units(self):
		if self.data.quantity == proto.Ingredient.VOLUME:
			return "ounces"
		elif self.data.quantity == proto.Ingredient.MASS:
			return "ounces"
		elif self.data.quantity == proto.Ingredient.COUNT:
			return ""
		elif self.data.quantity == proto.Ingredient.UNKNOWN:
			return self.data.quantity_type
		
	def __str__(self):
		if self.amount() % 1 == 0:
			return "[name: %s, quantity: %d %s]" % (self.name(), self.amount(), self.units())
		else:
			return "[name: %s, quantity: %.1f %s]" % (self.name(), self.amount(), self.units())
		
class Recipe(object):
	def __init__(self, id, title, prep_time, cook_time, ready_time, serving_size, ingredients):
		self.data = proto.Recipe()
		self.data.id = id
		self.data.title = title
		self.data.time.prep = prep_time
		self.data.time.cook = cook_time
		self.data.time.ready = ready_time
		self.data.servings = serving_size
		self.data.ingredients.extend( [ingr.data for ingr in ingredients] )

	def __str__(self):
		rstr = [
			"+-------------------------------------------------------+",
			"|    ** %s [[# %d]]** " % (self.data.title, self.data.id),
			"|    [Prep: %d min; Cook: %d min; Ready: %d min]" % (self.data.time.prep, self.data.time.cook, self.data.time.ready),
			"|    Servings: %d" % self.data.servings,
			"|    Contains: "
		]
		
		for ingr in self.data.ingredients:
			ingredient = Ingredient(ingr)
			rstr.append( '|      ' + str(ingredient) )

		rstr.append( "+-------------------------------------------------------+" )
		return '\n'.join(rstr)
