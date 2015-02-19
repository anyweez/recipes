package main

import (
	"fmt"
	"frontend/handlers"
	"frontend/state"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"lib/config"
	log "logging"
	"net/http"
)

type FrontendServer struct {
	Config config.RecipesConfig
	Router *mux.Router
}

// TODO: This should be the function that produces handlers.
type RouteHandler func(http.ResponseWriter, *http.Request)

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
var routes map[string]Registry

/**
 * Set the values of the route mapping.
 */
func init() {
	routes = map[string]Registry{
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

/**
 * This function produces handler functions that can be used to route http requests.
 * The returned functions use the mapping declared above in the `routes` variable
 * to invoke the desired function when the associated path is called. Gorilla is used
 * to route so URL params, etc are supported and passed on to the handler function.
 */
func ProduceHandler(path string, ss *state.SharedState) RouteHandler {

	return func(w http.ResponseWriter, r *http.Request) {
		// Log the web request event.
		le := log.New("web_request", log.Fields{
			"handled_path": path,
			"method":       r.Method,
		})

		// Check to ensure that a handler exists and retrieve it if it does.
		fn, exists := routes[path][r.Method]

		// Check to make sure that the handler exists. If so, run it.
		if exists {
			fn(w, r, ss, le)
			le.Update(log.STATUS_COMPLETE, "", nil)

			//			return nil
			// If no handler exists, log the issue and return it.
		} else {
			msg := fmt.Sprintf("No handler specified for method %s on path %s", r.Method, path)

			le.Update(log.STATUS_ERROR, msg, nil)
			//			return errors.New(msg)
		}
	}
}

/**
 * Create a new FrontendServer with the configuration parameters specified.
 * The configuration params will also be cascaded down to handlers, etc.
 */
func NewFrontendServer(conf config.RecipesConfig, le log.LogEvent) (FrontendServer, error) {
	// Save the configuration
	fes := FrontendServer{}
	fes.Config = conf

	// Load information for
	ss, err := state.NewSharedState(conf)

	if err != nil {
		le.Update(log.STATUS_ERROR, "Couldn't initialize shared state: "+err.Error(), nil)
		return fes, err
	}

	fes.Router = mux.NewRouter()

	/**
	 * Initialize all of the function handlers to handle the desired paths.
	 */
	for path := range routes {
		fes.Router.HandleFunc(path, ProduceHandler(path, ss))
	}

	/**
	 * Additional handlers to pass clientside assets (html, css, js, etc)
	 */
	// No-op handler for favicon.ico, since it'll otherwise generate an extra call to index.
	fes.Router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	// Serve any files in static/ directly from the filesystem.
	fes.Router.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		le := log.New("web_request", log.Fields{
			"handler": "<inline>",
			"path":    r.URL.Path[1:],
		})

		http.ServeFile(w, r, "html/"+r.URL.Path[1:])
		le.Update(log.STATUS_COMPLETE, "", nil)
	})

	le.Update(log.STATUS_OK, fmt.Sprintf("Frontend server initialized on port %d", conf.Frontend.Port), nil)
	return fes, nil
}

/**
 * Starts handling HTTP requests and continues handling them until the
 * process is stopped or an error occurs. If an error occurs that error
 * is returned back to the calling function.
 */
func (fes *FrontendServer) Start() error {
	http.Handle("/", fes.Router)

	return http.ListenAndServe(
		fmt.Sprintf(":%d", fes.Config.Frontend.Port),
		context.ClearHandler(http.DefaultServeMux),
	)
}
