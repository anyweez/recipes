package main

import (
	"flag"
	"lib/config"
	log "logging"
)

var CONFIG_LOCATION = flag.String("config", "recipes.conf", "The path to the configuration file.")

func main() {
	flag.Parse()
	le := log.New("frontend", nil)

	conf, err := config.New(*CONFIG_LOCATION)

	if err != nil {
		le := log.New("init", nil)
		le.Update(log.STATUS_FATAL, err.Error(), nil)
	}

	fes := NewFrontendServer(conf, le)
	// Start responding to requests. This is not expected to ever stop except
	// when explicitly killed.
	err = fes.Start()

	le.Update(log.STATUS_FATAL, "Couldn't listen on port 8088:"+err.Error(), nil)
}
