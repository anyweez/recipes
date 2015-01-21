package main

import (
	"flag"
	"fmt"
	"lib/config"
	log "logging"
	"net"
	"net/http"
	"net/rpc"
)

var PORT = flag.Int("port", 14501, "The port the process should listen on.")
var OUTPUT_QUADS = flag.String("out", "localhost:64210", "The file where the quads file should be output.")

var conf config.RecipesConfig

func main() {
	conf = config.New("recipes.conf")
	le := log.New("retriever", nil)
	
	retriever := new(Retriever)
	rpc.Register(retriever)
	rpc.HandleHTTP()

	// Start listening!
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))

	if err != nil {
		le.Update(log.STATUS_FATAL, "Couldn't start listening:" + err.Error(), nil)
	}

	le.Update(log.STATUS_OK, "Setup complete. Listening for RPC's on HTTP interface.", nil)
	http.Serve(l, nil)
}
