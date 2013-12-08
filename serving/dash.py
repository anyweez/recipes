
## Dash is a backend datastore service that retrieves responses for graph-based
## queries.
##
##
## Port range: (12100, 12149)
import proto.Recipes_pb2 as proto
import juggle.lib.juggle as juggle
from py2neo import neo4j

db = neo4j.GraphDatabaseService('http://localhost:7474/db/data')
ingredients = db.get_index(neo4j.Node, "Ingredients")

# Decoder converts request parameters into a series of queries that need
# to be executed.
def decode_request(request):
	print 'Received request.'
	
	req = proto.RecipeRequest()
	req.ParseFromString(request)
	
	return req

# Handler executes the queries that the decoder provides.
def handle(request):
	print 'Received request with %d specified ingredient(s).' % len(request.ingrids)
	
	recipes = {} # mapping from rid => # of ingredients this recipe appeared for
	for ingrid in request.ingrids:
		print 'finding recipes containing %s...' % ingrid
		ingr = ingredients.get('id', ingrid)
		
		if len(ingr) > 0:
			rels = db.match(start_node=ingr[0], rel_type="CONTAINS", bidirectional=True)
		
			for rel in rels:
				if recipes.has_key(rel.start_node['id']):
					recipes[rel.start_node['id']] += 1
				else:
					recipes[rel.start_node['id']] = 1
			
	# Only return the ID's of recipes where all ingredients appear.
	valid_recipes = [rid for rid, count in recipes.iteritems() if count == len(request.ingrids)]
		
	return valid_recipes

# Encoder packs the recipe ID's returned from the handler into a response
# proto.
def encode_response(response):
	print 'Returning response.'
	
	resp = proto.RecipeResponse()
	resp.test = "hello from dash! %d recipes found." % len(response)
	resp.recipe_ids.extend([str(rid) for rid in response])

	return resp.SerializeToString()

def run():
	service = juggle.Service('dash', 12100)
	service.handler(handle)
	service.encoder(encode_response)
	service.decoder(decode_request)

	service.start()

if __name__ == '__main__':
	print 'Starting dash.'
	run()
