
## Recipes web server. Doesn't do a whole lot; just serves the initial
## app and responds to the API requests defined at
## https://docs.google.com/document/d/15OcHFpliGclzjt9soDq2NiaBUfelvMyKt-1CzhT4hjU/edit?usp=sharing

import web, json, time, urllib
import juggle.lib.juggle as juggle
import proto.Recipes_pb2 as proto

from protobuf_to_dict import protobuf_to_dict

urls = (
	'/api/recipes', 'RecipeSearch',
	'/api/ingredients', 'Ingredients',
	'/', 'Index'
)

def encode_request(ingredient_list):
	rr = proto.RecipeRequest()
	rr.ingrids.extend(ingredient_list)

	return rr.SerializeToString()
	
def decode_response(response):
	rr = proto.RecipeResponse()
	rr.ParseFromString(response)
	
	return rr

class RecipeSearch(object):
	def GET(self):
		params = web.input()
		
		try:
			contains = urllib.unquote(params.contains).decode('utf-8')
			ingredient_list = contains.split(',')

			return self.find_recipes(ingredient_list)
		except AttributeError:
			# General recipe search.
			return json.dumps( {'error': 'Not supported yet!'} )

	# This method identifies recipes that use the ingredients specified
	# in 'ingredient_list' and returns information about those recipes.
	def find_recipes(self, ingredient_list):
		service = juggle.ServiceAPI('localhost', 12098)
		service.encoder(encode_request)
		service.decoder(decode_response)
		
		start_time = time.time()
		resp = service.query(ingredient_list)
		resp.duration = (time.time() - start_time) * 1000
		
		return json.dumps( self.prepare_response(resp) )

	# TODO: most of this will go away once the API is cleaned up (speed hacking now)
	def prepare_response(self, response):
		resp_dict = protobuf_to_dict(response)
		# Remove the recipe_ids field
		resp_dict.pop('recipe_ids', None)
		resp_dict.pop('test', None)
		resp_dict.pop('duration', None)
		
		if not resp_dict.has_key('recipes'):
			return {'data': [], 'debug': {'duration': response.duration}}
		
		for i, recipe in enumerate(resp_dict['recipes']):
			resp_dict['recipes'][i]['cook_time'] = recipe['time']['cook']
			resp_dict['recipes'][i]['prep_time'] = recipe['time']['prep']
			resp_dict['recipes'][i]['ready_time'] = recipe['time']['ready']
			
			for j, ingredient in enumerate(recipe['ingredients']):
				resp_dict['recipes'][i]['ingredients'][j]['units'] = 'ounce'
				resp_dict['recipes'][i]['ingredients'][j]['picture'] = 'not-implemented-yet.png'
		
				# Replace the ingrids repeated field with a single ID.
				if ingredient.has_key('ingrids'):
					resp_dict['recipes'][i]['ingredients'][j]['id'] = ingredient['ingrids'][0] 
					resp_dict.pop('ingrids', None)
				else:
					resp_dict['recipes'][i]['ingredients'][j]['id'] = -1
		
		return {'data': resp_dict['recipes'], 
		         'debug': {'duration': response.duration}}
		

class Index(object):
	def GET(self):
		with open('web/index.html') as fp:
			return fp.read()

class Ingredients(object):
	def GET(self):
		# also 'amount' (float) and 'amountUnits' (string)
		return json.dumps({ 'data': 
			[{'name': 'Cranberry sauce', 'id': 'm/0709jg'},
			 {'name': 'Chocolate', 'id': 'm/020vl'}]
		})

if __name__ == '__main__':
	app = web.application(urls, globals())
	app.run()
