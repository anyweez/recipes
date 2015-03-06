package main

import (
	"fmt"
	"frontend/handlers"
	"net/http"
	//	proto "proto"
	"testing"
)

const BASE_URL = "http://localhost:13033"

/**
 * Test to ensure that no groups are returned for a user that doesn't have any groups.
 * (also tests to make sure the login works correctly.)
 *
 */
func TestNoGroups(t *testing.T) {
	// Create the client.
	client, err := NewTestClient()

	if err != nil {
		t.Error(err.Error())
	}

	client.QueueRequest(TestLogin(handlers.LoginRequest{
		EmailAddress: "theo@bald.com",
	}))

	_, status, rerr := client.Run()

	if status != http.StatusOK {
		t.Error(fmt.Sprintf("Returned status code %d", status))
	}

	if rerr != nil {
		t.Error(rerr.Error())
	}
}
