package main

import (
	"frontend/handlers"
)

/**
 * Registries contain mappings between HTTP methods (GET, POST, etc) and
 * the handlers that should be used to fulfill the request.
 */
type Registry map[string]handlers.Handler

/**
 * Mapping of URL paths to core API handlers by method. The full specification
 * is available at https://github.com/luke-segars/dinder-docs.
 *
 * The values are actually set in init() in this file.
 */
var Routes map[string]Registry

/**
 * Set the values of the route mapping.
 */
func init() {
	Routes = map[string]Registry{
		"/api/users": Registry{
			"POST": handlers.CreateNewUser,
		},
		"/api/users/me": Registry{
			"GET": handlers.GetUser,
		},
		"/api/users/login": Registry{
			"POST": handlers.Login,
		},
		"/api/groups": Registry{
			"GET":  handlers.GetGroups,
			"POST": handlers.CreateGroup,
		},
		"/api/groups/{group_id}/user": Registry{
			"POST": handlers.AddUserToGroup,
		},
		"/api/meals/today": Registry{
			"GET": handlers.GetTodaysMeal,
		},
		"/api/meals/today/status": Registry{
			"POST": handlers.SetMealStatus,
		},
		"/api/meals/vote": Registry{
			"POST": handlers.SetMealVote,
		},
		"/api/recipes": Registry{
			"GET": handlers.GetBestRecipes,
		},
	}
}
