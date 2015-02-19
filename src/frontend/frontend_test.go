package main

import (
	"testing"
)

/**
 * Test to make sure that you can't request a list of groups for a user
 * unless they're logged in. This should return a 401 error as specified
 * in the documentation.
 */
func TestLoggedOutGroupsRequest(t *testing.T) {
	_, err := NewHttpTest("TestLoggedOutGroupsRequest", "GET", "/api/groups", 401)

	if err != nil {
		t.Errorf(err.Error())
	}
}
