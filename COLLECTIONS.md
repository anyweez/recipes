Collections
============

This document outlines all of the collections that are used for this application, how information
is stored, what it's added and removed by, and other useful information to understand the relevant
data flows. All collections are a part of the `recipes` database by default. In addition to the
key mentioned below, all tables are also keyed by a standard object ID (`_id`), though this interface
is rarely made available in API's, etc.

Note that the status below isn't representative of today's state but is representative of the
intended purpose of each collection.

Responses (default: `responses`)
---------------------------------

The responses table stores recent user responses (yes / no) for a given recipe in the context
of a group. This means that a user can provide different answers to the same recipe in different
groups (or at least the backend supports it; clients can keep all groups in sync if they desire).
Responses can be used to populate frontend UI elements and determine when a consensus has been
reached among everyone in the group. It is most accurately thought of as data for the current
session (per group); data from prior decisions is removed from this table and stored in longer
term but less readily-accessible logs.

Key: group_id (uint64)

Value: proto.RecipeResponses

Lifetime: until consensus is reached (per group); periodic cleanup as well.

Add:

	- `retrieve.PostRecipeResponses` RPC records a new response for a group member

Modify:

	- `retrieve.PostRecipeResponses` RPC will add a new response to an existing collection but will not replace existing responses
	  
Delete:

	- Automatically once consensus is reached
	- Periodically if `RecipeResponses` object is more than 24 hours old

Recipes (default: `recipes`)
------------------------------

The recipes table stores a context-free structured version of all recipes known to the system.
All recipe data (ingredients, name, image URL's, etc) are available in this table but no serving
state is maintained here. This table is the single source of truth for recipe data. All changes
to this table are made via an offline process so it can generally be assumed to be contention-free.

Key: recipe_id (string)

Value: proto.Recipe

Lifetime: forever

Add:

	- The `extract recipes` command reads unstructured HTML from the raw HTML table (see below), parses it, and inserts the structured output into this table.

Modify:

	- None.
	
Delete:

	- None.

Recommended recipes (default: `recommended`)
---------------------------------------------

The recommended recipes table copies a subset of the `recipes` table into a group-specific environment. Each group has a set of recipes that are "recommended"
based on an offline analysis pipeline and cached her for serving purposes. ServingStatus is maintained for recipes in this table.

Key: group_id (uint64)

Value: proto.RecommendedRecipes

Lifetime: periodic (roughly daily)

Add:

	- An offline recommendation pipeline generates a list of recipes per group and adds them to this table.
	
Modify:

	- `retrieve.GetBestRecipes` will modify the ServingStatus of this table. Recipes are not added or removed.
	
Delete:

	- Periodically replaced as new recommendations become available.

Ingredients (default: `ingredients`)
-------------------------------------

The ingredient table is the single source of truth for ingredient (component of a recipe) data. This table is context-free as well, meaning that no
quantities are stored, only fields that are generally true for any recipe that contains them.

Key: none

Value: proto.Ingredient

Lifetime: forever

Add:

	- The `extract ingredients` command parses a Freebase dump and extracts a set of entities that often serve as ingredients in recipes. This process stores it's output in this table.

Modify:

	- None.
	
Delete:

	- None.

Users (default: `users`)
-------------------------

The users table is the canonical source for user-related profile information.
It is rarely changed as part of an online process.

Key: user_id (uint64)

Value: proto.User

Lifetime: forever

Add: 

	- account creation
Modify: 

	- account modification (changing username, etc)
Delete: 

	- none

Groups (default: `groups`)
------------------------------

The groups table is the canonical source for information about groups. 

Key: group_id (uint64)

Value: proto.Group

Lifetime: forever

Add:

	- Creating a new account generates a single-user group containing the new user.
	- User-invoked action to create a new group.

Modify:

	- Adding / removing users from a group.
	
Delete:

	- User-invoked action to delete a group.

Raw HMTL (default: `scraper`)
------------------------------

This table contains the output of the web scraper and is not required to be accessible
for any serving-related processes. It contains raw HTML that can be parsed by the `extract`
process.

Key: none

Value: <TBD>

Lifetime: forever

Add:

	- Web scraper reads page content and stores it.
	
Modify:

	- None.
	
Delete:

	- None.
