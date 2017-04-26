package main

import (
	"log"
	"net/http"

	"io/ioutil"

	"fmt"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

// This will be use by our handler at /graphql
var schema *graphql.Schema

// This function runs at the start of the program
func init() {

	// We get the schema from the file, rather than having the schema inline here
	// I think will lead to better organizaiton of our own code
	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		// We will panic if we don't find the schema.graphql file in our server
		panic(err)
	}

	// We will use graphql-go library to parse our schema from "schema.graphql"
	// and the resolver is our struct that should fullfill everything in the Query
	// from our schema
	schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
	if err != nil {
		panic(err)
	}
}

func main() {
	// We will start a small server that reads our "graphiql.html" file and
	// responds with it, so we are able to have our own graphiql
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile("graphiql.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(page)
	}))

	// This is where our graphql server is handled, we declare "/graphql" as the route
	// where all our graphql requests will be directed to
	http.Handle("/graphql", &relay.Handler{Schema: schema})

	// We start the server by using ListenAndServe and we log if we have any error, hope not!
	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Resolver struct is the main resolver  wich we will use to fullfill
// queries and mutations that our schema.graphql defines
type Resolver struct{}

// Hello function resolves to hello: String! in the Query object in our schema.graphql
func (r *Resolver) Hello() string {
	return "world"
}
