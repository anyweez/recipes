package main

/**
 * The request object that contains all of the fields required to make a
 * best recipe RPC call.
 *
 * All fields are currently required.
 */
type BestRecipesRequest struct {
	// Will eventually be removed from the API. Currently useful for testing
	Seed    int64
	UserId  uint64
	GroupId uint64

	// The number of recipes desired, if possible (not guaranteed)
	Count int
}

type RecipeResponse struct {
	RecipeId string
	UserId   uint64
	GroupId  uint64
	// Yes if they accepted, no if they declined.
	Response bool
}

type RecipeResponseRequest struct {
	RecipeId string
	GroupId  uint64
}
