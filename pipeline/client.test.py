import juggle.lib.juggle as jugglelib
import logging
import recipes
import proto.Recipes_pb2 as proto

# Configure logging.
logging.basicConfig()
logging.getLogger().setLevel(logging.INFO)

# Sending a string, so this is a passthrough.
def encode(obj):
	return obj

# Will receive an Ingredient proto back, need to decode.
def decode(string):
	ingr = proto.Ingredient()
	ingr.ParseFromString(string)
	
	return recipes.Ingredient(ingr)

if __name__ == '__main__':
	# Set up the ingredient service to parse ingredient strings.
	ingredient_service = jugglelib.ServiceAPI('localhost', 19001)	
	ingredient_service.encoder(encode)
	ingredient_service.decoder(decode)
	
	tests = [
		"1 egg",
		"1 (15 ounce) can tomato sauce",
		"1/4 cup water",
		"1 envelope taco seasoning mix",
		"1 1/2 tablespoons chili powder",
		"1 tablespoon vegetable oil",
		"1 pound chicken breast tenderloins",
		"1 (15 ounce) can black beans, drained",
		"1/4 cup cream cheese",
		"1 cup shredded Mexican-style cheese blend, or more to taste",
		"1 (7.5 ounce) package corn bread mix",
		"1/3 cup milk",
		"1 (15 ounce) can pumpkin puree",
		"4 eggs",
		"1 cup vegetable oil",
		"2/3 cup water",
		"3 cups white sugar",
		"3 1/2 cups all-purpose flour",
		"2 teaspoons baking soda",
		"1 1/2 teaspoons salt",
		"1 teaspoon ground cinnamon",
		"1 teaspoon ground nutmeg",
		"1/2 teaspoon ground cloves",
		"1/4 teaspoon ground ginger",
		"1 pound boneless pork loin chops, pounded thin"
	]

	for test in tests:
		print ingredient_service.query(test)
