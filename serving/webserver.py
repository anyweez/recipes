
## Recipes web server. Doesn't do a whole lot; just serves the initial
## app and responds to the API requests defined at
## https://docs.google.com/document/d/15OcHFpliGclzjt9soDq2NiaBUfelvMyKt-1CzhT4hjU/edit?usp=sharing

import web, json, time
import juggle.lib.juggle as juggle
import proto.Recipes_pb2 as proto

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
		
		ingredient_list = params.contains.split(',')

		service = juggle.ServiceAPI('localhost', 12098)
		service.encoder(encode_request)
		service.decoder(decode_response)
		
		start_time = time.time()
		resp = service.query(ingredient_list)
		end_time = time.time() - start_time

		return 'Response: %s [%.1f ms]' % (resp.test, end_time * 1000)

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
