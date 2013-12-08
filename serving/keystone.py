
## Keystone is the middleware the functions between the web server
## frontend and the databases. Its general tasks include:
##   - fetching documents
##   - ranking
##
## It should be possible to receive requests from multiple frontends
## simultaneously and to load balance over multiple backend DB's as well.
##
## Port range: (12050, 12099)

import juggle.lib.juggle as juggle
import proto.Recipes_pb2 as proto

def frontend_request(request):
	return request

## Serialize the RecipeResponse to send to the webserver.
def frontend_response(response):
	return response.SerializeToString()

def backend_request(request):
	return request

# Deserialize the response so keystone can use it to fetch additional
# recipe information.
def backend_response(response):
	resp = proto.RecipeResponse()
	resp.ParseFromString(response)
	
	return resp

def handle(request):
	print 'initializing connection to backend'
	backend = juggle.ServiceAPI('localhost', 12100)
	backend.encoder(backend_request)
	backend.decoder(backend_response)
	
	# Response containing a bunch of recipe ID's. Fetch all of the recipes
	# and add the additional content.
	resp = backend.query(request)

	# TODO: fetch recipe docs from mongo
	
	return resp

def run():
	frontend_service = juggle.Service('keystone', 12098)

	# Handler receives a requests from the frontend and passes the
	# request through to a backend. It then retrieves information about
	# the 
	frontend_service.handler(handle)
	
	# Encoder should handle the retrieval of the recipe ID's that are
	# returned from the backend.
	frontend_service.encoder(frontend_response)
	# Decoder is a passthrough.
	frontend_service.decoder(frontend_request)
	
	print 'Awaiting connections.'
	frontend_service.start()
	
if __name__ == '__main__':
	run()
