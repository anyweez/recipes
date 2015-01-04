package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var PORT = flag.Int("port", 14501, "The port the process should listen on.")
var OUTPUT_QUADS = flag.String("out", "localhost:64210", "The file where the quads file should be output.")

func main() {
	retriever := new(Retriever)
	rpc.Register(retriever)
	rpc.HandleHTTP()
	log.Println("Launching retriever...")

	// Start listening!
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *PORT))

	if err != nil {
		log.Fatal("Couldn't start listening:" + err.Error())
	}

	log.Println("Setup complete. Listening for RPC's on HTTP interface.")
	http.Serve(l, nil)
}
