package main

/**
 * This file contains a series of functions that all return http.Request objects that simulate
 * what real client requests look like. These functions can be added together in series to
 * produce different test cases.
 */

import (
	"bytes"
	"encoding/json"
	"frontend/handlers"
	"net/http"
)

const BASE_URL = "http://localhost:13033"

func TestLogin(user handlers.LoginRequest) (*http.Request, error) {
	data, jerr := json.Marshal(user)
	if jerr != nil {
		return nil, jerr
	}

	// Form the request and fire it off.
	req, rerr := http.NewRequest("POST", BASE_URL+"/api/users/login", bytes.NewReader(data))
	if rerr != nil {
		return nil, rerr
	}

	return req, nil
}

/*
func TestGetGroups() (*http.Request, error) {

}
*/
