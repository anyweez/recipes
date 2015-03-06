package main

import (
	"bytes"
	"encoding/json"
	"frontend/handlers"
	"net/http"
)

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
