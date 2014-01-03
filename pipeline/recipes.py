import proto.Recipes_pb2 as proto

## Contains data structures used for storing, printing, and manipulating
## recipe-related data.

class Ingredient(object):
	def __init__(self, ingr):
		self.data = ingr
		
	def amount(self):
		return self.data.quantity.original_amount
	
	def name(self):
		return self.data.name
	
	def units(self):
		if self.data.quantity.original_type == proto.Ingredient.CUP:
			return 'cup' if self.data.quantity.original_amount < 1.05 else 'cups'
		elif self.data.quantity.original_type == proto.Ingredient.TABLESPOON:
			return 'tablespoon' if self.data.quantity.original_amount < 1.05 else 'tablespoons'
		elif self.data.quantity.original_type == proto.Ingredient.TEASPOON:
			return 'teaspoon' if self.data.quantity.original_amount < 1.05 else 'teaspoons'
		elif self.data.quantity.original_type == proto.Ingredient.QUART:
			return 'quart' if self.data.quantity.original_amount < 1.05 else 'quarts'
		elif self.data.quantity.original_type == proto.Ingredient.GALLON:
			return 'gallon' if self.data.quantity.original_amount < 1.05 else 'gallons'
		elif self.data.quantity.original_type == proto.Ingredient.POUND:
			return 'pound' if self.data.quantity.original_amount < 1.05 else 'pounds'
		elif self.data.quantity.original_type == proto.Ingredient.PINT:
			return 'pint' if self.data.quantity.original_amount < 1.05 else 'pints'
		elif self.data.quantity.original_type == proto.Ingredient.CUSTOM:
			return self.data.quantity.original_unit_string
		
	def __str__(self):
		if self.amount() % 1 == 0:
			return "[name: %s, quantity: %d %s, custom type? %r]" % (self.name(), self.amount(), self.units(), self.data.quantity.original_type == proto.Ingredient.CUSTOM)
		else:
			return "[name: %s, quantity: %.1f %s, custom type? %r]" % (self.name(), self.amount(), self.units(), self.data.quantity.original_type == proto.Ingredient.CUSTOM)
		
class Recipe(object):
	def __init__(self, id, name, prep_time, cook_time, ready_time, serving_size, ingredients):
		self.data = proto.Recipe()
		self.data.id = id
		self.data.name = name
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
