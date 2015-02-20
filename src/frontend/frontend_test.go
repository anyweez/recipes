package main

import (
	"testing"
)

/**
 * Test to make sure that all expected endpoints exist and nothing tragic
 * happens when an empty body is sent to them.
 */
func TestEndpointsExist(t *testing.T) {
	for path := range Routes {
		registry := Routes[path]

		// Check each method (there should only be one)
		// function per method but there can be several
		// methods.
		for method := range registry {
			// Note that this is inverted from the normal use csae; checking
			// to see if 404 is returned and if so that's a problem.
			_, err := NewHttpTest(path, method, path, 404)

			//
			if err == nil {
				t.Errorf("Endpoint `%s` doesn't exist.", path)
			}
		}
	}
}

/**
 * Test to make sure that you can't successfully get anything back but a
 * 401 error for any of the API calls that require authentication.
 *
 * This test should test every publicly exposed method except for methods
 * that don't require the user to be logged in (login methods, for example).
 */
func TestLoggedOutRequests(t *testing.T) {
	get := map[string]string{
		"TestLoggedOutGroupsRequest": "/api/groups",
		"TestLoggedOutMeals":         "/api/meals/today",
		"TestLoggedOutUserSelf":      "/api/users/me",
		"TestLoggedOutRecipes":       "/api/recipes",
	}

	// Test all get methods.
	for name := range get {
		_, err := NewHttpTest(name, "GET", get[name], 401)

		if err != nil {
			t.Errorf("%s: %s", name, err.Error())
		}
	}

	post := map[string]string{
		"TestCreateGroup":    "/api/groups",
		"TestAddUserToGroup": "/api/groups/0/user",
		"TestSetMealStatus":  "/api/meals/today/status",
		"TestSetMealVote":    "/api/meals/vote",
	}

	// Test all get methods.
	for name := range post {
		_, err := NewHttpTest(name, "POST", post[name], 401)

		if err != nil {
			t.Errorf("%s: %s", name, err.Error())
		}
	}
}
