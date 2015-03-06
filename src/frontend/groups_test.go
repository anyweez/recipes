package main

import (
	"fmt"
	"frontend/handlers"
	"net/http"
	//	proto "proto"
	"testing"
)

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

	// Login
	_, lstatus, lerr := client.DoTest(TestLogin(handlers.LoginRequest{
		EmailAddress: "theo@bald.com",
	}))

	if lerr != nil {
		t.Error(lerr.Error())
	}

	if lstatus != 200 {
		t.Error("Issues with logging in.")
	}

	// Retrieve
	_, gstatus, gerr := client.DoTest(TestGetGroups())

	if gerr != nil {
		t.Error(gerr.Error())
	}

	if gstate != 200 {
		t.Error("Error retrieving groups.")
	}

	_, status, rerr := client.Run()

	if status != http.StatusOK {
		t.Error(fmt.Sprintf("Returned status code %d", status))
	}

	if rerr != nil {
		t.Error(rerr.Error())
	}
}
