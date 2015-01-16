recipes
=======

The goal of this project is to provide an ingestion pipeline and serving infrastructure
for an ingredient-based recipe search engine. The major serving components are designed
to be more generic so that they can be used for other low latency, scalable QPS, retrieval-
based projects in the future.

Usage
-----
Running the build file will install all required dependencies and build the key binaries for
both the offline pipeline and the serving infrastructure. (Almost) all of the tools access
configuration parameters stores in recipes.conf so make sure these are representative of
your environment.

Dependencies
------------
*Cayley* Cayley is a graph database that's used for recipe/ingredient analysis and more complex
retrieval tasks.

*MongoDB* Mongo is used as the primary datastore for structured (non-graph) data. Most of
the binaries associated with this project access a shared MongoDB database.
