package main

import (
	"errors"
	"fmt"
	"lib/config"
	log "logging"
	"net/http"
	"net/http/httptest"
)

const (
	ConfigPath = "/home/luke/git/recipes/recipes.conf"
)

/**
 * Run a test by sending an HTTP request to `path` using the specified `method.` This function
 * also checks that the response code that comes back is equivalent to `expected` and will return
 * an error if not.
 */
func NewHttpTest(name string, method string, path string, expected int) (*FrontendServer, error) {
	// Initialize the test configuration.
	le := log.New(name, nil)
	conf, cerr := config.New(ConfigPath)

	if cerr != nil {
		return &FrontendServer{}, cerr
	}

	fes, ferr := NewFrontendServer(conf, le)

	if ferr != nil {
		return &FrontendServer{}, cerr
	}

	// Once we've got the frontend server initialized, the work is the same as if the FES had already
	// existed.
	return NextHttpTest(name, method, path, expected, &fes)
}

func NextHttpTest(name string, method string, path string, expected int, fes *FrontendServer) (*FrontendServer, error) {
	// Make the request and record the response.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return fes, err
	}

	w := httptest.NewRecorder()
	fes.Router.ServeHTTP(w, req)

	// Check whether the response code matches the expected value. If not, we
	// should return an error.
	if w.Code != expected {
		return fes, errors.New(fmt.Sprintf("Returned HTTP code %d instead of expected %d", w.Code, expected))
	}

	return fes, nil
}
