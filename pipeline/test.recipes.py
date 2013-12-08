import sys
import proto.Recipes_pb2 as proto

## Prints out information about recipes in a serialized RecipeBook.

with open(sys.argv[1]) as fp:
	book = proto.RecipeBook()
	book.ParseFromString(fp.read())
	
	for recipe in book.recipes:
		print recipe.id
